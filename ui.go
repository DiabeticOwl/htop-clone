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

	cpuProgresses    []progress.Model
	memoryProgresses []progress.Model

	cpuTable       table.Model
	memoryTable    table.Model
	disksTable     table.Model
	processesTable table.Model

	// Window's width.
	Width  int
	Height int
}

func NewModel() model {
	vMemoryInfo, sMemoryInfo := extractMemoryInfo()

	teaModel := model{
		CpuInfo:     extractCpuInfo(),
		Processes:   extractProcessesInfo(),
		VMemoryInfo: vMemoryInfo,
		SMemoryInfo: sMemoryInfo,
	}

	opts := []progress.Option{
		progress.WithDefaultGradient(),
	}
	for range teaModel.CpuInfo {
		pBar := progress.New(opts...)
		pBar.PercentFormat = " %.2f%%"

		teaModel.cpuProgresses = append(teaModel.cpuProgresses, pBar)
	}

	opts = append(opts, progress.WithoutPercentage())
	teaModel.memoryProgresses = []progress.Model{
		progress.New(opts...), // One for each type of memory.
		progress.New(opts...),
	}

	pCount := int(math.Ceil(float64(len(teaModel.Processes)) / 20))
	teaModel.processesTable = newProcessesTable(teaModel, pCount)

	teaModel.cpuTable = newCpuTable(teaModel)
	teaModel.memoryTable = newMemoryTable(teaModel)
	teaModel.disksTable = newDisksTable(teaModel)

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
		for i := range m.cpuProgresses {
			m.cpuProgresses[i].Width = int(float64(msg.Width) * 0.15)
		}

		for i := range m.memoryProgresses {
			m.memoryProgresses[i].Width = int(float64(msg.Width) * 0.20)
		}

		// -2 as an arbitrary margin.
		m.Width = msg.Width - 2

		// -33 as an experimental value for calculating the
		// amount of processes per page.
		pCount := msg.Height - 33
		if pCount <= 0 {
			pCount = 2
		}

		m.cpuTable = newCpuTable(m)
		if msg.Height <= 20 {
			m.processesTable = table.New([]table.Column{})
			m.disksTable = table.New([]table.Column{})
			m.memoryTable = table.New([]table.Column{})
		} else {
			m.processesTable = newProcessesTable(m, pCount)
			m.disksTable = newDisksTable(m)
			m.memoryTable = newMemoryTable(m)
		}

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
	s := lipgloss.NewStyle().Padding(1).Render(m.cpuTable.View())
	s += lipgloss.NewStyle().Padding(1).Render(m.memoryTable.View())
	s += lipgloss.NewStyle().Padding(1).Render(m.disksTable.View())
	s += lipgloss.NewStyle().Padding(1).Render(m.processesTable.View())

	// Send the UI for rendering.
	return s
}
