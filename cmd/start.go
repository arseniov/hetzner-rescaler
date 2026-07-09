package cmd

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/broadcast"
	"github.com/jonamat/hetzner-rescaler/internal/crypto"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/scheduler"
	"github.com/jonamat/hetzner-rescaler/internal/store"
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

	keyring, err := crypto.NewKeyringFromBytes(key)
	if err != nil {
		fmt.Println("keyring:", err)
		return
	}
	keyringHolder = keyring

	// Lifecycle hub: receive server create/update/delete events so the
	// scheduler can keep its per-server goroutines in sync without polling.
	lifecycleHub := broadcast.NewHub[store.ServerLifecycleEvent]()
	st.SetServerLifecycleHub(lifecycleHub)

	// Scheduler: per-server goroutines that drive scheduled / auto_promote
	// rescales. The apiResolve uses the same per-project factory the
	// API server uses, so multi-project deployments get full coverage
	// (no phase-1 warning anymore).
	factory := apiFactory(st)
	sched := scheduler.New(st, func(_ context.Context, projectID int64) (hetzner.API, error) {
		return factory(projectID)
	}, scheduler.RealClock{}, 30*time.Second)

	// Lifecycle: bind to signal context and run Attach in a goroutine.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := sched.Attach(ctx, lifecycleHub); err != nil && !errors.Is(err, context.Canceled) {
			fmt.Println("scheduler attach:", err)
		}
	}()

	fmt.Printf("Scheduler started for %d existing servers. Press Ctrl+C to stop.\n", len(servers))

	// Block until SIGINT/SIGTERM, then drain.
	<-ctx.Done()
	fmt.Println("\nShutting down...")
	sched.Stop()
}
