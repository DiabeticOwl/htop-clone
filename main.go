// htop-clone runs a bubbletea application that displays the user's computer
// main health statistics.
package main

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	teaModel := model{
		CpuInfo: extractCpuInfo(interval),
	}
	for range teaModel.CpuInfo {
		opts := []progress.Option{
			progress.WithDefaultGradient(),
		}
		pBar := progress.New(opts...)
		pBar.PercentFormat = " %.2f%%"

		teaModel.progresses = append(teaModel.progresses, pBar)
	}

	p := tea.NewProgram(teaModel)
	if err := p.Start(); err != nil {
		panic(err)
	}
}
