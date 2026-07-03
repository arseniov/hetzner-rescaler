package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

/* Serve command (phase 1 stub) */
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the HTTP API server (phase 2)",
	Long:  "Phase 1 stub: prints a message and exits. The HTTP server ships in phase 2 (see docs/superpowers/specs/2026-06-29-hetzner-rescaler-web-design.md).",
	Run:   runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) {
	fmt.Println("serve: not yet implemented (phase 2). The CLI commands config, start, try, status, and migrate are functional in phase 1.")
}