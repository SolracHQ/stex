package main

import (
	"fmt"
	"os"

	"github.com/SolracHQ/stex/tui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
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

	program := tea.NewProgram(tui.New(path))
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
