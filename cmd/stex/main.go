// stex is an interactive TUI to explore where disk space is being used. The program is a
// single binary that takes a path argument (defaults to the current directory) and starts a
// Bubble Tea program.
package main

import (
	"fmt"
	"os"

	"github.com/SolracHQ/stex/internal/app"
	"github.com/SolracHQ/stex/internal/config"

	flag "github.com/spf13/pflag"

	tea "charm.land/bubbletea/v2"
)

func main() {
	icons := flag.BoolP("icons", "i", false, "start with emoji icons enabled")
	showAll := flag.BoolP("show-all", "a", false, "start with hidden files shown")
	noLive := flag.BoolP("no-live-filter", "L", false, "disable live filter (compile on enter)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: stex [flags] [path]\n\nFlags:\n")
		fmt.Fprintf(os.Stderr, "  -i, --icons          start with emoji icons enabled\n")
		fmt.Fprintf(os.Stderr, "  -a, --show-all       start with hidden files shown\n")
		fmt.Fprintf(os.Stderr, "  -L, --no-live-filter disable live filter (compile on enter)\n")
		fmt.Fprintf(os.Stderr, "  -h, --help           show this help\n")
	}
	flag.Parse()

	path := "."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: %s is a file, not a directory\n", path)
		os.Exit(1)
	}

	cfg := config.DefaultConfig()
	if *icons {
		cfg.ShowIcons = true
	}
	if *showAll {
		cfg.ShowHidden = true
	}
	if *noLive {
		cfg.LiveFilter = false
	} else {
		cfg.LiveFilter = true
	}

	program := tea.NewProgram(app.New(path, cfg))
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
