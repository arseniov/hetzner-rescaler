package cmd

import (
	"fmt"
	"io"

	"github.com/jonamat/hetzner-rescaler/internal/store"
	"github.com/spf13/cobra"
)

/* Status command */
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print all configured projects, servers, and recent events",
	Run:   runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) {
	st, err := openStore()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer st.Close()
	if err := printStatus(st, cmd.OutOrStdout()); err != nil {
		fmt.Println("error:", err)
	}
}

func printStatus(st *store.Store, w io.Writer) error {
	projects, err := st.ListProjects()
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		fmt.Fprintln(w, "No projects configured. Run `hetzner-rescaler config` to add one.")
		return nil
	}
	for _, p := range projects {
		fmt.Fprintf(w, "Project %d: %s\n", p.ID, p.Name)
		servers, err := st.ListServersByProject(p.ID)
		if err != nil {
			return err
		}
		for _, srv := range servers {
			fmt.Fprintf(w, "  Server %d: %s (hcloud_id=%d, mode=%s, base=%s, top=%s, fallback=%v, tz=%s)\n",
				srv.ID, srv.Name, srv.HCloudServerID, srv.Mode,
				srv.BaseServerType, srv.TopServerType, srv.FallbackChain, srv.Timezone)
			events, err := st.ListEventsByServer(srv.ID, 5)
			if err != nil {
				return err
			}
			for _, e := range events {
				fmt.Fprintf(w, "    [%s] %s: %s -> %s ok=%v\n",
					e.StartedAt.Format("2006-01-02 15:04"), e.Kind, e.FromType, e.ToType, e.OK)
			}
		}
	}
	return nil
}
