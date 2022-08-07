// htop-clone runs a bubbletea application that displays the user's computer
// main health statistics.
package main

import (
	"flag"
	"os"
	"runtime/pprof"

	tea "github.com/charmbracelet/bubbletea"
)

// Tag for activating Go's profiling tools.
// The given value will be used as the name of the file in which the profile
// will be saved.
var cpuProfile = flag.String("cpuProfile", "", "Writes CPU profile to a file.")

func main() {
	// Parses the flags.
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			panic(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		panic(err)
	}
}
