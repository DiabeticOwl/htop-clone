package main

import (
	"math"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	interval time.Duration = time.Second
)

type model struct {
	CpuInfo []float64
	// Virtual Memory.
	VMemoryInfo map[string]interface{}
	// Swap Memory.
	SMemoryInfo map[string]interface{}
	Processes   []map[string]interface{}
	DisksInfo   []map[string]interface{}

	progresses []progress.Model
	msgWidth   int

	cpuTable       table.Model
	memoryTable    table.Model
	disksTable     table.Model
	processesTable table.Model
}

func NewModel() model {
	teaModel := model{
		CpuInfo:     extractCpuInfo(),
		Processes:   extractProcessesInfo(),
		cpuTable:    newCpuTable(),
		memoryTable: newMemoryTable(),
		disksTable:  newDisksTable(),
	}
	for range teaModel.CpuInfo {
		opts := []progress.Option{
			progress.WithDefaultGradient(),
		}
		pBar := progress.New(opts...)
		pBar.PercentFormat = " %.2f%%"

		teaModel.progresses = append(teaModel.progresses, pBar)
	}

	pCount := int(math.Ceil(float64(len(teaModel.Processes)) / 20))
	teaModel.processesTable = newProcessesTable(pCount)

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

	m.memoryTable, cmd = m.memoryTable.Update(msg)
	cmds = append(cmds, cmd)

	m.disksTable, cmd = m.disksTable.Update(msg)
	cmds = append(cmds, cmd)

	m.processesTable, cmd = m.processesTable.Update(msg)
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

		m.cpuTable = m.cpuTable.WithRows(generateCpuTableRows(m))
		m.memoryTable = m.memoryTable.WithRows(generateMemoryTableRows(m))
		m.disksTable = m.disksTable.WithRows(generateDisksTableRows(m))
		m.processesTable = m.processesTable.WithRows(generateProcessesTableRows(m))

		return m, tea.Batch(cmds...)
	}

	// Return the updated model and no command is left to run.
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := lipgloss.NewStyle().Padding(1).Render(m.cpuTable.View() + "\n")
	s += lipgloss.NewStyle().Padding(1).Render(m.memoryTable.View() + "\n")
	s += lipgloss.NewStyle().Padding(1).Render(m.disksTable.View() + "\n")
	s += lipgloss.NewStyle().Padding(1).Render(m.processesTable.View() + "\n")

	// Applications's footer.
	s += "\nPress q or Ctrl+C to exit.\n"

	// Send the UI for rendering.
	return s
}
