package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
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

	progresses []progress.Model
	msgWidth   int
}

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (_ model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if k := msg.String(); k == "q" || k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		for i := range m.progresses {
			m.progresses[i].Width = int(float64(msg.Width) * 0.15)
		}

		m.msgWidth = msg.Width

		return m, nil

	case tickMsg:
		m.CpuInfo = extractCpuInfo(interval)
		m.VMemoryInfo, m.SMemoryInfo = extractMemoryInfo()
		m.DisksInfo = extractDiskInfo()
		m.Processes = extractProcessesInfo()

		cmds := []tea.Cmd{
			tick(),
		}

		for i := range m.progresses {
			cmds = append(cmds, m.progresses[i].SetPercent(m.CpuInfo[i]/100))
		}

		sort.Sort(sort.Reverse(byCpuUsage(m.Processes)))

		// Extract first 20 processes.
		m.Processes = m.Processes[:20]

		return m, tea.Batch(cmds...)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		var cmds []tea.Cmd

		for i := range m.progresses {
			progressModel, cmd := m.progresses[i].Update(msg)
			m.progresses[i] = progressModel.(progress.Model)

			cmds = append(cmds, cmd)
		}

		return m, tea.Batch(cmds...)
	}

	// Return the updated model and no command is left to run.
	return m, nil
}

func (m model) View() string {
	s := cpuTableView(m)

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
	s += "* PID \t| User \t\t| Priority | Niceness | CPU Usage Percentage | Executable Path *\n"
	for _, p := range m.Processes {
		sP := "| %d  | %s | %d \t| %d \t| %.4f%% | %s"

		s += fmt.Sprintf(sP, p.PId, p.User, p.Priority, p.Niceness,
			p.CpuPercentage, p.ExeP)
		s += "\n"
	}
	s += "\n************ Processes info. ************\n"

	// Applications's footer.
	s += "\nPress q or Ctrl+C to exit.\n"

	// Send the UI for rendering.
	return s
}

func cpuTableView(m model) string {
	// The amount of rows will be the amount of cpu divided by 3.
	rowCount := int(math.Ceil(float64(len(m.CpuInfo)) / 3))

	// Title.
	t := " CPU Usage Percentage "
	tPad := strings.Repeat("-", int(float64(m.msgWidth)*0.35))

	sTable := fmt.Sprintf("•%s%s%s•\n", tPad, t, tPad)

	// A counter is used to assign each value of the m.CpuInfo slice
	// to its respective index in the cpuTable matrix.
	var counter int
	for c := 0; c < 3; c++ {
		for r := 0; r < rowCount; r++ {
			valStr := " CPU #%d: %s\t"
			sTable += fmt.Sprintf(valStr, counter, m.progresses[counter].View())

			counter++

			if r == rowCount-1 {
				sTable += "\n"
			}
		}
	}

	tPad = strings.Repeat("-", int(float64(m.msgWidth)*0.35)*2+len(t))
	sTable += fmt.Sprintf("•%s•\n", tPad)

	return sTable
}
