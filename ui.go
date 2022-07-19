package main

import (
	"fmt"
	"math"
	"strings"
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
	// // Application's header.
	// s := "\n************ Cpu info. ************\n"
	// s += "| "
	// sCpu := "| "

	// for i := 0; i < len(m.CpuInfo); i++ {
	// 	s += fmt.Sprintf("CPU #%d Usage | ", i)
	// 	sCpu += fmt.Sprintf("\t %.2f%%|", m.CpuInfo[i])
	// }

	// s += "\n"
	// s += sCpu + "\n"
	// s += "\n************ Cpu info. ************\n"
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
	// Width of each column. Establishing 3 columns per row.
	columnLengths := make([]int, 3)
	// The amount of rows will be the amount of cpu divided by 3.
	rowCount := int(math.Ceil(float64(len(m.CpuInfo)) / 3))

	cpuTable := make([][]string, rowCount)
	for i := 0; i < rowCount; i++ {
		cpuTable[i] = make([]string, 3)
	}

	// A counter is used to assign each value of the m.CpuInfo slice
	// to its respective index in the cpuTable matrix.
	var counter int
	for c := range columnLengths {
		for r := 0; r < rowCount; r++ {
			valStr := "CPU #%d: %.2f%%"
			cpuTable[r][c] = fmt.Sprintf(valStr, counter, m.CpuInfo[counter])

			counter++
		}
	}

	// Establishing the largest width in each column.
	for _, r := range cpuTable {
		for i, val := range r {
			if len(val) > columnLengths[i] {
				columnLengths[i] = len(val)
			}
		}
	}

	var rowLength int
	for _, c := range columnLengths {
		rowLength += c + 3 // + 3 For extra padding in each value.
	}
	rowLength += 1 // +1 For the last "|" in each row.

	// If there is nothing to show the default value of rowLength will
	// be columnCount ^ 2 + 1.
	if rowLength == 10 {
		return ""
	}

	// Title.
	t := " CPU Usage Percentage "
	tLen := len(t)
	tPad := strings.Repeat("-", (rowLength-2)/2-tLen/2)

	sTable := fmt.Sprintf("•%s%s%s•\n", tPad, t, tPad)

	for i, row := range cpuTable {
		for j, val := range row {
			// Formats each row with the corresponding width and value.
			sTable += fmt.Sprintf("| %-*s ", columnLengths[j], val)

			// If j corresponds to the last column, add the last
			// decoration to the row.
			if j == len(row)-1 {
				sTable += "|\n"
			}
		}

		// After formatting the header and all rows, formats the footer.
		if i == len(cpuTable)-1 {
			sTable += fmt.Sprintf("•%s•\n", strings.Repeat("-", rowLength-2))
		}
	}

	return sTable
}
