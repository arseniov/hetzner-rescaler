package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/crypto"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/scheduler"
	"github.com/spf13/cobra"
)

/* Start command */
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the scheduler loop (reads from SQLite)",
	Run:   runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) {
	st, err := openStore()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer st.Close()

	servers, err := st.ListAllServers()
	if err != nil {
		fmt.Println("list servers:", err)
		return
	}
	if len(servers) == 0 {
		fmt.Println("No servers configured. Run `hetzner-rescaler config` to add one.")
		return
	}

	key, err := loadOrGenerateKey()
	if err != nil {
		fmt.Println("key:", err)
		return
	}

	// Build one Hetzner API per project
	apiByProject := map[int64]hetzner.API{}
	for _, srv := range servers {
		if _, ok := apiByProject[srv.ProjectID]; ok {
			continue
		}
		proj, err := st.GetProject(srv.ProjectID)
		if err != nil {
			fmt.Println("project:", err)
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
		apiByProject[srv.ProjectID] = api
	}

	// Phase 1: a single scheduler instance with one project's API.
	// Multi-project in one process is phase 2; for now, run one `start` per
	// project if you have more than one.
	//
	// BUG FIX: do NOT index the map by `0` — SQLite AUTOINCREMENT IDs start
	// at 1, so apiByProject[0] would always be nil. Pick the first API we
	// actually built.
	var firstAPI hetzner.API
	var firstProjectID int64
	for id, api := range apiByProject {
		firstAPI = api
		firstProjectID = id
		break
	}
	if firstAPI == nil {
		fmt.Println("error: failed to build API client for any project")
		return
	}
	if len(apiByProject) > 1 {
		fmt.Printf("WARNING: multiple projects detected; phase 1 schedules only project %d. Run one `hetzner-rescaler start` per project for full coverage.\n", firstProjectID)
	}

	sched := scheduler.New(st, firstAPI, scheduler.RealClock{}, 30*time.Second)

	// Register only the servers that belong to the selected project.
	registered := 0
	for _, srv := range servers {
		if srv.ProjectID != firstProjectID {
			continue
		}
		sched.Add(srv.ID)
		registered++
	}

	fmt.Printf("Scheduler started for %d servers. Press Ctrl+C to stop.\n", registered)

	// Handle SIGINT/SIGTERM
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Println("\nShutting down...")
		sched.Stop()
	}()

	sched.Run()
}
