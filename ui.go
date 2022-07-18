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
	// Virtual Memory.
	VMemoryInfo virtualMemoryInfo
	// Swap Memory.
	SMemoryInfo swapMemoryInfo
	Processes   []processInfo
	DisksInfo   []diskInfo
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
		m.VMemoryInfo, m.SMemoryInfo = extractMemoryInfo()
		m.DisksInfo = extractDiskInfo()
		m.Processes = extractProcessesInfo()

		return m, tick()
	}

	// Return the updated model and no command is left to run.
	return m, nil
}

func (m model) View() string {
	// Application's header.
	s := "\n************ Cpu info. ************\n"
	s += "| "
	sCpu := "| "

	for i := 0; i < len(m.CpuInfo); i++ {
		s += fmt.Sprintf("CPU #%d Usage | ", i)
		sCpu += fmt.Sprintf("\t %.2f%%|", m.CpuInfo[i])
	}

	s += "\n"
	s += sCpu + "\n"
	s += "\n************ Cpu info. ************\n"

	s += "\n************ Memory info. ************\n"

	s += "\n------------ Virtual Memory. ------------\n"
	sVm := "| Total: %.2f GB | Used: %.2f GB | Available: %.2f GB | "
	sVm += "UsedPercent: %.4f%% |"
	s += fmt.Sprintf(sVm, m.VMemoryInfo.Total, m.VMemoryInfo.Used,
		m.VMemoryInfo.Available, m.VMemoryInfo.UsedPercent)
	s += "\n------------ Virtual Memory. ------------\n"

	s += "\n------------ Swap Memory. ------------\n"
	sSm := "| Total: %.2f GB | Used: %.2f GB | Free: %.2f GB | "
	sSm += "UsedPercent: %.4f%% |"
	s += fmt.Sprintf(sSm, m.SMemoryInfo.Total, m.SMemoryInfo.Used,
		m.SMemoryInfo.Free, m.SMemoryInfo.UsedPercent)
	s += "\n------------ Swap Memory. ------------\n"

	s += "\n************ Memory info. ************\n"

	s += "\n************ Disks info. ************\n"
	for _, d := range m.DisksInfo {
		sD := "| Device: %s | MountPath: %s | TotalSize: %.2f GB | FreeSize: "
		sD += "%.2f GB | UsedSize: %.2f GB |"

		s += fmt.Sprintf(sD, d.Device, d.MountPath, d.TotalSize, d.FreeSize,
			d.UsedSize)
		s += "\n"
	}
	s += "\n************ Disks info. ************\n"

	s += "\n************ Processes info. ************\n"
	for _, p := range m.Processes {
		sP := "| PID: %d | User: %s | Priority: %d | Niceness: %d | "
		sP += "CPU Percentage: %.4f%% | Command: %s"

		s += fmt.Sprintf(sP, p.PId, p.User, p.Priority, p.Niceness,
			p.CpuPercentage, p.Cmdline)
		s += "\n"
	}
	s += "\n************ Processes info. ************\n"

	// Applications's footer.
	s += "\nPress q or Ctrl+C to exit.\n"

	// Send the UI for rendering.
	return s
}
