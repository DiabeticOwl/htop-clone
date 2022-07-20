// htop-clone runs a bubbletea application that displays the user's computer
// main health statistics.
package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(NewModel())
	if err := p.Start(); err != nil {
		panic(err)
	}
}
