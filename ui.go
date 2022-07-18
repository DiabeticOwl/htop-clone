package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	interval time.Duration = time.Second
)

type model struct {
	CpuInfo []float64
	// Both Virtual and Swap Memory.
	MemoryInfo [][]float64
	Processes  []processInfo
	DisksInfo  []diskInfo
}

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if k := msg.String(); k == "q" || k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch msg.(type) {
	case tickMsg:
		m.CpuInfo = extractCpuInfo(interval)
		return m, tick()
	}

	// Return the updated model and no command is left to run.
	return m, nil
}

func (m model) View() string {
	// Application's header.
	s := "\n------------ Cpu info. ------------\n"
	s += "|"

	sCpu := "|"

	for i := 0; i < len(m.CpuInfo); i++ {
		s += fmt.Sprintf("\t CPU #%d Usage |", i)
		sCpu += fmt.Sprintf("\t %.2f%%", m.CpuInfo[i])
	}

	s += "\n"
	s += sCpu + "\n"

	// Applications's footer.
	s += "\nPress q or Ctrl+C to exit.\n"

	// Send the UI for rendering.
	return s
}
