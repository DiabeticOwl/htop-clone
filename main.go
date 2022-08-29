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
	// // Channel that will accept one signal.
	// signChan := make(chan os.Signal, 1)
	// // Sends a SIGTERM signal to the given channel when it is
	// // received by the application.
	// signal.Notify(signChan, syscall.SIGTERM)

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
	// go func() {
	// 	if sig := <-signChan; sig != nil {
	// 		p.Quit()
	// 	}
	// }()
	if err := p.Start(); err != nil {
		panic(err)
	}
}
