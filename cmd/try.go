package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/crypto"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/rescaler"
	"github.com/jonamat/hetzner-rescaler/internal/store"
	"github.com/spf13/cobra"
)

/* Try command */
var tryCmd = &cobra.Command{
	Use:   "try <server-id> <up|down>",
	Short: "One-shot rescale of a single server",
	Long:  "Rescales a single registered server up or down immediately. The server must be registered first (see `hetzner-rescaler config`).",
	Args:  cobra.ExactArgs(2),
	Run:   runTry,
}

func init() {
	rootCmd.AddCommand(tryCmd)
}

func runTry(cmd *cobra.Command, args []string) {
	serverID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fmt.Println("invalid server id:", args[0])
		return
	}
	direction := args[1]
	if direction != "up" && direction != "down" {
		fmt.Println("direction must be 'up' or 'down'")
		return
	}

	st, err := openStore()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer st.Close()

	srv, err := st.GetServer(serverID)
	if err != nil {
		fmt.Println("server not found:", err)
		return
	}
	proj, err := st.GetProject(srv.ProjectID)
	if err != nil {
		fmt.Println("project not found:", err)
		return
	}
	key, err := loadOrGenerateKey()
	if err != nil {
		fmt.Println("key:", err)
		return
	}
	token, err := crypto.Decrypt(key, proj.HCloudTokenEncrypted, proj.HCloudTokenNonce)
	if err != nil {
		fmt.Println("decrypt token:", err)
		return
	}
	api, err := hetzner.NewClient(string(token))
	if err != nil {
		fmt.Println("client:", err)
		return
	}

	hsrv, err := api.GetServer(context.Background(), srv.HCloudServerID)
	if err != nil {
		fmt.Println("get server:", err)
		return
	}
	target := srv.TopServerType
	if direction == "down" {
		target = srv.BaseServerType
	}
	if hsrv.ServerType.Name == target {
		fmt.Printf("server is already at %s, nothing to do\n", target)
		return
	}

	fmt.Printf("Rescaling server %d (%s) from %s to %s...\n", hsrv.ID, hsrv.Name, hsrv.ServerType.Name, target)
	start := time.Now().UTC()
	used, err := rescaler.RescaleWithFallback(context.Background(), api, hsrv, target, srv.FallbackChain)
	finished := time.Now().UTC()
	if err != nil {
		_, _ = st.AppendEvent(store.Event{
			ServerID:    srv.ID,
			Kind:        "rescale_failed",
			FromType:    hsrv.ServerType.Name,
			ToType:      target,
			StartedAt:   start,
			FinishedAt:  finished, // value type (not pointer) — Task 6 review change
			OK:          false,
			Error:       err.Error(),
			TriggeredBy: "manual",
		})
		fmt.Println("rescale failed:", err)
		return
	}
	kind := "rescale_up"
	if direction == "down" {
		kind = "rescale_down"
	}
	_, _ = st.AppendEvent(store.Event{
		ServerID:    srv.ID,
		Kind:        kind,
		FromType:    hsrv.ServerType.Name,
		ToType:      used,
		StartedAt:   start,
		FinishedAt:  finished, // value type
		OK:          true,
		TriggeredBy: "manual",
	})
	fmt.Printf("OK: rescaled to %s\n", used)
}
