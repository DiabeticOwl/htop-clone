// htop-clone runs a bubbletea application that displays the user's computer
// main health statistics.
package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(model{})
	if err := p.Start(); err != nil {
		panic(err)
	}
}
