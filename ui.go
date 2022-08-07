// File that describes the UI/UX that the user will interact with.
package main

import (
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
	VMemoryInfo memoryInfo
	// Swap Memory.
	SMemoryInfo memoryInfo
	Processes   []processInfo
	DisksInfo   []diskInfo

	cpuProgresses    []progress.Model
	memoryProgresses []progress.Model

	cpuTable       table.Model
	memoryTable    table.Model
	disksTable     table.Model
	processesTable table.Model

	// Window's width.
	Width int
	// Window's height.
	Height int
}

// NewModel initializes the model that BubbleTea will use.
func NewModel() model {
	teaModel := model{
		CpuInfo: extractCpuInfo(),
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

	return teaModel
}

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if k := msg.String(); k == "q" || k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		} else if k == "a" || k == "A" {
			m.disksTable = m.disksTable.PageUp()
			return m, nil
		} else if k == "d" || k == "D" {
			m.disksTable = m.disksTable.PageDown()
			return m, nil
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
			m.memoryProgresses[i].Width = int(float64(msg.Width) * 0.15)
		}

		// -2 as an arbitrary margin.
		m.Width = msg.Width - 2
		m.Height = msg.Height

		// -33 as an experimental value for calculating the amount of processes
		// per page.
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
		m.processesTable = m.processesTable.WithRows(generateProcessesTableRows(m))

		var pCount int
		if len(m.DisksInfo) > 2 {
			pCount = 2
		}
		m.disksTable = m.disksTable.WithRows(generateDisksTableRows(m)).WithPageSize(pCount)

		return m, tea.Batch(cmds...)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	// Render each component in the UI with its given style.
	s := lipgloss.NewStyle().Padding(0, 1, 1).Render(m.cpuTable.View())
	s += lipgloss.NewStyle().Padding(1).Render(m.memoryTable.View())
	s += lipgloss.NewStyle().Padding(1).Render(m.disksTable.View())
	s += lipgloss.NewStyle().Padding(1).Render(m.processesTable.View())
	s += "\n a/d for the disks table. ↑ / ↓ / ← / → for processes table."

	// Send the UI for rendering in the terminal screen.
	return s
}
