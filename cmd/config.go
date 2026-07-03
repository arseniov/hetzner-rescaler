package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/jonamat/hetzner-rescaler/internal/crypto"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/store"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

/* Config command */
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Interactively add or edit projects, servers, modes, and windows",
	Run:   runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) {
	st, err := openStore()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer st.Close()
	key, err := loadOrGenerateKey()
	if err != nil {
		fmt.Println("key:", err)
		return
	}

	for {
		choice, err := promptMainMenu(st)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		switch choice {
		case "add project":
			if err := addProjectFlow(st, key); err != nil {
				fmt.Println("error:", err)
			}
		case "add server":
			if err := addServerFlow(st, key); err != nil {
				fmt.Println("error:", err)
			}
		case "exit":
			return
		}
	}
}

func promptMainMenu(st *store.Store) (string, error) {
	items := []string{"add project", "add server", "exit"}
	sel := promptui.Select{Label: "What do you want to do?", Items: items}
	_, choice, err := sel.Run()
	return choice, err
}

func addProjectFlow(st *store.Store, key []byte) error {
	token, err := promptToken()
	if err != nil {
		return err
	}
	api, err := hetzner.NewClient(token)
	if err != nil {
		return err
	}
	name, err := promptProjectName()
	if err != nil {
		return err
	}
	enc, nonce, err := crypto.Encrypt(key, []byte(token))
	if err != nil {
		return err
	}
	p, err := st.CreateProject(name, enc, nonce)
	if err != nil {
		return err
	}
	color.Green("Created project %q (id=%d). Token validated by initial ping.", name, p.ID)
	// Sanity check the token by listing servers
	servers, err := api.ListServers(context.Background())
	if err != nil {
		color.Yellow("Warning: could not list servers with this token: %v", err)
	} else {
		color.Green("Token is valid; %d servers visible in this project.", len(servers))
	}
	return nil
}

func addServerFlow(st *store.Store, key []byte) error {
	projects, err := st.ListProjects()
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		return fmt.Errorf("no projects yet; add a project first")
	}
	projectItems := make([]string, len(projects))
	for i, p := range projects {
		projectItems[i] = fmt.Sprintf("%d: %s", p.ID, p.Name)
	}
	sel := promptui.Select{Label: "Which project?", Items: projectItems}
	idx, _, err := sel.Run()
	if err != nil {
		return err
	}
	proj := projects[idx]

	token, err := crypto.Decrypt(key, proj.HCloudTokenEncrypted, proj.HCloudTokenNonce)
	if err != nil {
		return err
	}
	api, err := hetzner.NewClient(string(token))
	if err != nil {
		return err
	}

	hcloudServers, err := api.ListServers(context.Background())
	if err != nil {
		return err
	}
	if len(hcloudServers) == 0 {
		return fmt.Errorf("no servers in this project")
	}
	names := make([]string, len(hcloudServers))
	for i, s := range hcloudServers {
		names[i] = fmt.Sprintf("%d: %s (%s)", s.ID, s.Name, s.ServerType.Name)
	}
	sel = promptui.Select{Label: "Which server?", Items: names}
	idx, _, err = sel.Run()
	if err != nil {
		return err
	}
	hsrv := hcloudServers[idx]

	types, err := api.ListServerTypes(context.Background())
	if err != nil {
		return err
	}
	typeNames := make([]string, len(types))
	for i, t := range types {
		typeNames[i] = t.Name
	}
	sel = promptui.Select{Label: "Base (down) type?", Items: typeNames}
	_, baseName, err := sel.Run()
	if err != nil {
		return err
	}
	sel = promptui.Select{Label: "Top (up) type?", Items: typeNames}
	_, topName, err := sel.Run()
	if err != nil {
		return err
	}
	if baseName == topName {
		return fmt.Errorf("base and top types must differ")
	}

	chainInput, err := promptChain(topName, baseName)
	if err != nil {
		return err
	}

	mode, err := promptMode()
	if err != nil {
		return err
	}
	tz, err := promptTimezone()
	if err != nil {
		return err
	}

	srv, err := st.CreateServer(proj.ID, store.Server{
		HCloudServerID: hsrv.ID,
		Name:           hsrv.Name,
		BaseServerType: baseName,
		TopServerType:  topName,
		FallbackChain:  chainInput,
		Mode:           mode,
		Timezone:       tz,
	})
	if err != nil {
		return err
	}
	color.Green("Registered server %q (id=%d, mode=%s)", hsrv.Name, srv.ID, mode)

	// If scheduled, immediately prompt for the first window
	if mode == "scheduled" {
		if err := addWindowFlow(st, srv); err != nil {
			return err
		}
	}
	return nil
}

func addWindowFlow(st *store.Store, srv *store.Server) error {
	label, err := promptString("Window label", "weekday 9-19")
	if err != nil {
		return err
	}
	dow, err := promptDaysOfWeek()
	if err != nil {
		return err
	}
	start, err := promptString("Start time (HH:MM)", "09:00")
	if err != nil {
		return err
	}
	stop, err := promptString("Stop time (HH:MM)", "19:00")
	if err != nil {
		return err
	}
	// Default the window's target type to the server's top type. (Phase 1
	// does not yet let windows target a different type than the server's
	// top; that's a phase 2 enhancement.)
	_, err = st.CreateWindow(srv.ID, store.Window{
		Label: label, DaysOfWeek: dow,
		StartTime: start, StopTime: stop,
		TargetType: srv.TopServerType,
		Enabled:    true,
	})
	return err
}

func promptToken() (string, error) {
	p := promptui.Prompt{Label: "Hetzner Cloud token", Mask: '*'}
	return p.Run()
}
func promptProjectName() (string, error) {
	return promptString("Project name", "")
}
func promptMode() (string, error) {
	items := []string{"scheduled", "auto_promote", "manual"}
	sel := promptui.Select{Label: "Rescale mode", Items: items}
	_, mode, err := sel.Run()
	return mode, err
}
func promptTimezone() (string, error) {
	return promptString("Timezone (IANA, e.g. Europe/Rome)", "UTC")
}
func promptChain(top, base string) ([]string, error) {
	color.Yellow("Fallback chain (in order). First entry should be '%s'; last should be '%s'.", top, base)
	color.Yellow("Comma-separated. Press enter for the default chain [%s, %s].", top, base)
	defaultChain := top + "," + base
	raw, err := promptString("Fallback chain", defaultChain)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("chain must have at least one entry")
	}
	if out[len(out)-1] != base {
		return nil, fmt.Errorf("chain must end with the base type %q; got %q", base, out[len(out)-1])
	}
	return out, nil
}

// promptDaysOfWeek uses bit positions matching Go's time.Weekday, where
// Sunday=0 ... Saturday=6.
func promptDaysOfWeek() (int, error) {
	items := []string{"Mon-Fri (weekdays)", "Mon-Sun (all days)", "Sat-Sun (weekend only)"}
	sel := promptui.Select{Label: "Days of week", Items: items}
	idx, _, err := sel.Run()
	if err != nil {
		return 0, err
	}
	switch idx {
	case 0:
		return 0b00111110, nil // Mon-Fri: bits 1..5
	case 1:
		return 0b01111111, nil // Mon-Sun: all 7 bits
	case 2:
		return 0b01000001, nil // Sat-Sun: bits 0,6
	}
	return 0, fmt.Errorf("unreachable")
}
func promptString(label, def string) (string, error) {
	p := promptui.Prompt{Label: label, Default: def}
	return p.Run()
}
func promptInt(label string, def int) (int, error) {
	raw, err := promptString(label, strconv.Itoa(def))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(raw)
}