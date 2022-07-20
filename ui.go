package main

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	interval time.Duration = time.Second

	cpuTableTitle           = "CPU Usage Percentage"
	cpuTableMaxColumnAmount = 3

	// https://pkg.go.dev/github.com/evertras/bubble-table@v0.14.4/table?utm_source=gopls#NewFlexColumn
	columnStandardFlexFactor = 1
	columnKeyCpuTable        = "cpuTable"
)

var (
	styleBase = (lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("#a7a")).
		BorderBackground(lipgloss.Color("#a38")).
		Align(lipgloss.Center))
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

	cpuTable table.Model
}

func newCpuTable() table.Model {
	columns := []table.Column{
		table.NewFlexColumn(columnKeyCpuTable, cpuTableTitle,
			columnStandardFlexFactor),
	}

	return (table.
		New(columns).
		BorderRounded().
		WithBaseStyle(styleBase).
		WithTargetWidth(180))
}

func (m model) generateCpuTableData() []table.Row {
	valStr := "CPU #%d: %s "
	rowCount := int(math.Ceil(float64(len(m.CpuInfo)) / cpuTableMaxColumnAmount))

	if rowCount*cpuTableMaxColumnAmount != len(m.CpuInfo) {
		s := fmt.Sprintf(`The amount of columns or rows is incorrect.
		The product of both should be %d.`, len(m.CpuInfo))

		panic(s)
	}

	var rows []table.Row
	for i := 0; i < rowCount; i++ {
		r := ""
		for c := 0; c < cpuTableMaxColumnAmount; c++ {
			if c > 0 {
				// index of the row + index of the column * the amount of values
				// in one column.
				index := i + c*rowCount
				r += fmt.Sprintf(valStr, index, m.progresses[index].View())
			} else {
				r += fmt.Sprintf(valStr, i, m.progresses[i].View())
			}
		}

		nRow := table.NewRow(table.RowData{
			columnKeyCpuTable: r,
		}).WithStyle(lipgloss.NewStyle())
		rows = append(rows, nRow)
	}

	return rows
}

func NewModel() model {
	teaModel := model{
		CpuInfo:  extractCpuInfo(),
		cpuTable: newCpuTable(),
	}
	for range teaModel.CpuInfo {
		opts := []progress.Option{
			progress.WithDefaultGradient(),
		}
		pBar := progress.New(opts...)
		pBar.PercentFormat = " %.2f%%"

		teaModel.progresses = append(teaModel.progresses, pBar)
	}

	return teaModel
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

	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.cpuTable, cmd = m.cpuTable.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		for i := range m.progresses {
			m.progresses[i].Width = int(float64(msg.Width) * 0.15)
		}

		m.msgWidth = msg.Width

		return m, tea.Batch(cmds...)

	case tickMsg:
		m.CpuInfo = extractCpuInfo()
		m.VMemoryInfo, m.SMemoryInfo = extractMemoryInfo()
		m.DisksInfo = extractDiskInfo()
		m.Processes = extractProcessesInfo()

		cmds = append(cmds, tick())

		for i := range m.progresses {
			cmds = append(cmds, m.progresses[i].SetPercent(m.CpuInfo[i]/100))
		}

		sort.Sort(sort.Reverse(byCpuUsage(m.Processes)))

		// Extract first 20 processes.
		m.Processes = m.Processes[:20]

		m.cpuTable = m.cpuTable.WithRows(m.generateCpuTableData())

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
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := lipgloss.NewStyle().Padding(1).Render(m.cpuTable.View() + "\n")
	// s := lipgloss.NewStyle().Render(cpuTableView)
	// s := cpuTableView(m)

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
