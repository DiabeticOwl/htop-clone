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

	// Minimum terminal's window height for showing one tables.
	minimumHeightOneTable = 9
	// Minimum terminal's window height for showing two tables.
	minimumHeightTwoTables = 14
	// Minimum terminal's window height for showing three tables.
	minimumHeightThreeTables = 24
	// Minimum terminal's window height for showing all tables.
	minimumHeightAllTables = 33
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
	// Initial model instance with CpuInfo filled.
	teaModel := model{
		CpuInfo: extractCpuInfo(),
	}

	// Creating progress bars for the Cpu and Memory tables.
	opts := []progress.Option{
		progress.WithDefaultGradient(),
	}
	for range teaModel.CpuInfo {
		pBar := progress.New(opts...)
		pBar.PercentFormat = "%05.2f%% "
		pBar.PercentageRightPosition = false

		teaModel.cpuProgresses = append(teaModel.cpuProgresses, pBar)
	}

	opts = append(opts, progress.WithoutPercentage())
	teaModel.memoryProgresses = []progress.Model{
		progress.New(opts...), // One for each type of memory.
		progress.New(opts...),
	}

	return teaModel
}

// Time type for the tick function of BubbleTea.
type tickMsg time.Time

// tick returns a signal or "tick" after a set interval of time passes.
func tick() tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init initializes the program's model and returns its first "tick".
// It is attached to the program's model as it is required for the BubbleTea package.
func (model) Init() tea.Cmd {
	return tick()
}

// Update takes a given message and acts on the model according to what type or
// value it has. It also send the message to each table in order to update it
// and append the table's response to an array of commands to run at the end of
// this function. Update returns the updated model and a command or a batch of them.
// In this program there are three expected types of messages:
// * Keyword presses.
// * Terminal's window resizing.
// * "Tick"s, where another one is returned and hence creating a loop.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Given a keyword pressed return the updated model and a command.
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

	// Update each table and hold an array of commands to run.
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
		// -2 as an arbitrary margin.
		m.Width = msg.Width - 2
		m.Height = msg.Height

		if msg.Height < minimumHeightOneTable {
			m.cpuTable = table.New([]table.Column{})
			m.memoryTable = table.New([]table.Column{})
			m.disksTable = table.New([]table.Column{})
			m.processesTable = table.New([]table.Column{})

			return m, tea.Batch(cmds...)
		}

		// Progress bar width.
		pWidth := int(float64(msg.Width) * 0.15)

		for i := range m.cpuProgresses {
			m.cpuProgresses[i].Width = pWidth
		}

		for i := range m.memoryProgresses {
			m.memoryProgresses[i].Width = pWidth
		}

		// cpuTable will always be shown.
		m.cpuTable = newCpuTable(m)
		switch h := m.Height; {
		case h >= minimumHeightTwoTables && h < minimumHeightThreeTables:
			m.memoryTable = newMemoryTable(m)
		case h >= minimumHeightThreeTables && h < minimumHeightAllTables:
			m.memoryTable = newMemoryTable(m)
			m.disksTable = newDisksTable(m)
		case h >= minimumHeightAllTables:
			pCount := msg.Height - minimumHeightAllTables
			if pCount <= 0 {
				pCount = 2
			}

			m.memoryTable = newMemoryTable(m)
			m.disksTable = newDisksTable(m)
			m.processesTable = newProcessesTable(m, pCount)
		}

		return m, tea.Batch(cmds...)

	// Update each table each "tick".
	case tickMsg:
		m.CpuInfo = extractCpuInfo()
		m.VMemoryInfo, m.SMemoryInfo = extractMemoryInfo()
		m.DisksInfo = extractDiskInfo()
		m.Processes = extractProcessesInfo()

		m.cpuTable = m.cpuTable.WithRows(generateCpuTableRows(m))
		m.memoryTable = m.memoryTable.WithRows(generateMemoryTableRows(m))
		m.processesTable = m.processesTable.WithRows(generateProcessesTableRows(m))

		var pCount int
		if len(m.DisksInfo) > 2 {
			pCount = 2
		}
		m.disksTable = m.disksTable.WithRows(generateDisksTableRows(m)).WithPageSize(pCount)

		// Added "tick" to the final array of commands to run.
		cmds = append(cmds, tick())

		return m, tea.Batch(cmds...)
	}

	return m, tea.Batch(cmds...)
}

// View will render the current program state from a returned string.
func (m model) View() string {
	var s string

	if m.Height < minimumHeightOneTable {
		s = "\nWindow size is too small to show something."
	} else {
		// Render each component in the UI with its given style.
		s = lipgloss.NewStyle().Padding(0, 1, 1).Render(m.cpuTable.View())
		switch h := m.Height; {
		case h >= minimumHeightTwoTables && h < minimumHeightThreeTables:
			s += lipgloss.NewStyle().Padding(1).Render(m.memoryTable.View())
		case h >= minimumHeightThreeTables && h < minimumHeightAllTables:
			s += lipgloss.NewStyle().Padding(1).Render(m.memoryTable.View())
			s += lipgloss.NewStyle().Padding(1).Render(m.disksTable.View())
			s += "\n a/d for the disks table navigation."
		case h >= minimumHeightAllTables:
			s += lipgloss.NewStyle().Padding(1).Render(m.memoryTable.View())
			s += lipgloss.NewStyle().Padding(1).Render(m.disksTable.View())
			s += lipgloss.NewStyle().Padding(1).Render(m.processesTable.View())
			s += "\n a/d for the disks table, ↑ / ↓ / ← / → for processes table navigation."
		}
	}

	return s
}
