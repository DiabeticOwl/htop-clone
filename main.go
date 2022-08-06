// htop-clone runs a bubbletea application that displays the user's computer
// main health statistics.
package main

import (
	"flag"
	"os"
	"runtime/pprof"

	tea "github.com/charmbracelet/bubbletea"
)

var cpuProfile = flag.String("cpuProfile", "", "Writes CPU profile to a file.")

func main() {
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
