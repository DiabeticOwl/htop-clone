// htop-clone runs a bubbletea application that displays the user's computer
// main health metrics.
package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		// Many unaccounted errors can come from sys calls.
		// They are unlikely to occur.
		panic(err)
	}
}
