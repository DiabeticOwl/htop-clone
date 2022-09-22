package main

import (
	"fmt"
	"math"
	"runtime"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	cpuTableTitle = "CPU Usage Percentage"
	// Amount of columns in one row.
	cpuTableMaxColumnAmount = 4

	columnKeyCpuTable = "cpuTable"

	columnKeyVirtualMemory      = "virtualMemory"
	columnKeyVirtualMemoryTitle = "Virtual Memory"
	columnKeySwapMemory         = "swapMemory"
	columnKeySwapMemoryTitle    = "Swap Memory"
)

const (
	// https://pkg.go.dev/github.com/evertras/bubble-table/table?utm_source=gopls#NewFlexColumn
	columnDefaultFlexFactor = 1
	columnLargerFlexFactor  = iota * 2
	columnHugeFlexFactor
	columnLargestFlexFactor
)

var (
	styleBase = (lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("#92DCE5")).
			Align(lipgloss.Center))

	standardRowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#CEEFF3"))
)

// newCpuTable will instantiate the CPU information table with its assigned
// columns. This should be ran only when the application starts or it resizes.
func newCpuTable(m model) table.Model {
	columns := []table.Column{
		table.NewFlexColumn(columnKeyCpuTable, cpuTableTitle,
			columnDefaultFlexFactor),
	}

	return table.
		New(columns).
		BorderRounded().
		WithBaseStyle(styleBase).
		WithTargetWidth(m.Width)
}

// generateCpuTableRows will generate all the rows that will be rendered into
// the CPU information table. This should be ran each time the application updates.
func generateCpuTableRows(m model) []table.Row {
	cpuFmt := "CPU #%02d:"
	// The number of cores divided by the number of chosen columns will be
	// ceiled in order to calculate a fixed number of rows in which display
	// all the CPU info in an uniform manner.
	rowCount := int(math.Ceil(float64(len(m.CpuInfo)) / cpuTableMaxColumnAmount))

	// When the number of cpuTableMaxColumnAmount doesn't satisfy the equation
	// len(m.CpuInfo) / rowCount = cpuTableMaxColumnAmount, the program can't
	// display all cores properly.
	if rowCount*cpuTableMaxColumnAmount != len(m.CpuInfo) {
		s := "The amount of columns or rows is incorrect.\n"
		s = fmt.Sprintf("%sThe product of both should be %d.", s, len(m.CpuInfo))

		panic(s)
	}

	// Since each bubble-table row needs to be returned with all the information
	// once we need to populate each with the information of the corresponding
	// columns. For example, using a 12 core CPU and 4 columns each row:
	// Row 0 : "CPU #00 m.cpuProgresses[0] CPU #03 m.cpuProgresses[3] CPU #06 m.cpuProgresses[6] CPU #09 m.cpuProgresses[9]"
	// Row 1 : "CPU #01 m.cpuProgresses[1] CPU #04 m.cpuProgresses[4] CPU #07 m.cpuProgresses[7] CPU #10 m.cpuProgresses[10]"
	// Row 2 : "CPU #02 m.cpuProgresses[2] CPU #05 m.cpuProgresses[5] CPU #08 m.cpuProgresses[8] CPU #11 m.cpuProgresses[11]"
	var rows []table.Row
	for i := 0; i < rowCount; i++ {
		r := ""
		for c := 0; c < cpuTableMaxColumnAmount; c++ {
			if c > 0 {
				// index of the row + index of the column * the amount of values
				// in one column.
				index := i + c*rowCount

				r += fmt.Sprintf("%s %s", standardRowStyle.SetString(fmt.Sprintf(cpuFmt, index)).String(), m.cpuProgresses[index].ViewAs(m.CpuInfo[i]/100))
			} else {
				r += fmt.Sprintf("%s %s", standardRowStyle.SetString(fmt.Sprintf(cpuFmt, i)).String(), m.cpuProgresses[i].ViewAs(m.CpuInfo[i]/100))
			}
		}

		nRow := table.NewRow(table.RowData{
			columnKeyCpuTable: r,
		})
		rows = append(rows, nRow)
	}

	return rows
}

// newMemoryTable will instantiate the RAM information table with its assigned
// columns. This should be ran only when the application starts or it resizes.
func newMemoryTable(m model) table.Model {
	columns := []table.Column{
		table.NewFlexColumn(columnKeyVirtualMemory, columnKeyVirtualMemoryTitle,
			columnDefaultFlexFactor),
		table.NewFlexColumn(columnKeySwapMemory, columnKeySwapMemoryTitle,
			columnDefaultFlexFactor),
	}

	return table.
		New(columns).
		BorderRounded().
		WithBaseStyle(styleBase).
		WithTargetWidth(m.Width)
}

// generateMemoryTableRows will generate all the rows that will be rendered into
// the RAM information table. This should be ran each time the application updates.
func generateMemoryTableRows(m model) []table.Row {
	vMemoryProg := m.memoryProgresses[0].ViewAs(m.VMemoryInfo.UsedPercent / 100)
	vMemoryView := fmt.Sprintf("%s %s", standardRowStyle.SetString(fmt.Sprintf("%.2f GB/%.2f GB", m.VMemoryInfo.Used, m.VMemoryInfo.Total)).String(), vMemoryProg)

	sMemoryProg := m.memoryProgresses[1].ViewAs(m.SMemoryInfo.UsedPercent / 100)
	sMemoryView := fmt.Sprintf("%s %s", standardRowStyle.SetString(fmt.Sprintf("%.2f GB/%.2f GB", m.SMemoryInfo.Used, m.SMemoryInfo.Total)).String(), sMemoryProg)

	rows := []table.Row{
		table.NewRow(table.RowData{
			columnKeyVirtualMemory: vMemoryView,
			columnKeySwapMemory:    sMemoryView,
		}),
	}

	return rows
}

// newDisksTable will instantiate the disks information table with its assigned
// columns. This should be ran only when the application starts or it resizes.
func newDisksTable(m model) table.Model {
	fsTypeCol := table.NewFlexColumn("FsType", "File System Type", columnDefaultFlexFactor)
	deviceCol := table.NewFlexColumn("Device", "Device", columnDefaultFlexFactor)
	mountPathCol := table.NewFlexColumn("MountPath", "Mount Path", columnHugeFlexFactor)
	totalSizeCol := table.NewFlexColumn("TotalSize", "Total Size", columnDefaultFlexFactor)
	freeSizeCol := table.NewFlexColumn("FreeSize", "Free Size", columnDefaultFlexFactor)
	usedSizeCol := table.NewFlexColumn("UsedSize", "Used Size", columnDefaultFlexFactor)

	columns := []table.Column{fsTypeCol, deviceCol, mountPathCol, totalSizeCol, freeSizeCol, usedSizeCol}

	return table.
		New(columns).
		BorderRounded().
		WithBaseStyle(styleBase.Copy().Align(lipgloss.Left)).
		WithTargetWidth(m.Width).
		SortByAsc("FsType").
		ThenSortByAsc("MountPath")
}

// generateDisksTableRows will generate all the rows that will be rendered into
// the disks information table. This should be ran each time the application updates.
func generateDisksTableRows(m model) []table.Row {
	var rows []table.Row

	for _, disk := range m.DisksInfo {
		rowData := make(table.RowData)

		rowData["FsType"] = disk.FsType
		rowData["Device"] = disk.Device
		rowData["MountPath"] = disk.MountPath
		rowData["TotalSize"] = fmt.Sprintf("%2.f GB", disk.TotalSize)
		rowData["FreeSize"] = fmt.Sprintf("%2.f GB", disk.FreeSize)
		rowData["UsedSize"] = fmt.Sprintf("%2.f GB", disk.UsedSize)

		row := table.NewRow(rowData).WithStyle(standardRowStyle)
		rows = append(rows, row)
	}

	return rows
}

// newProcessesTable will instantiate the processes information table with its
// assigned columns. This should be ran only when the application starts or it
// resizes.
func newProcessesTable(m model, pCount int) table.Model {
	pIdCol := table.NewFlexColumn("PId", "Process ID", columnDefaultFlexFactor)
	prioCol := table.NewFlexColumn("Priority", "Priority", columnDefaultFlexFactor)

	uCol := table.NewFlexColumn("User", "Username", columnLargerFlexFactor)
	cPcgCol := table.NewFlexColumn("CpuPercentage", "CPU Usage Percentage", columnLargerFlexFactor).WithFormatString("%.1f%%")
	nCol := table.NewFlexColumn("Name", "Name", columnLargerFlexFactor)

	exePCol := table.NewFlexColumn("ExeP", "Executable Path", columnHugeFlexFactor)

	cmdlineCol := table.NewFlexColumn("Cmdline", "Command", columnLargestFlexFactor)

	columns := []table.Column{pIdCol, prioCol, uCol, cPcgCol, nCol, exePCol, cmdlineCol}

	// Not showing the exePCol as name and executable path are the same in darwin
	// based systems.
	if runtime.GOOS == "darwin" {
		columns = []table.Column{pIdCol, prioCol, uCol, cPcgCol, nCol, cmdlineCol}
	}

	return table.
		New(columns).
		BorderRounded().
		WithBaseStyle(styleBase.Copy().Align(lipgloss.Left)).
		WithTargetWidth(m.Width).
		WithPageSize(pCount).
		SortByDesc("CpuPercentage").
		Focused(true)
}

// generateProcessesTableRows will generate all the rows that will be rendered
// into the processes information table. This should be ran each time the
// application updates.
func generateProcessesTableRows(m model) []table.Row {
	var rows []table.Row

	for _, process := range m.Processes {
		rowData := make(table.RowData)

		rowData["PId"] = process.PId
		rowData["Priority"] = process.Priority
		rowData["User"] = process.User
		rowData["CpuPercentage"] = process.CpuPercentage
		rowData["Name"] = process.Name
		rowData["ExeP"] = process.ExeP
		rowData["Cmdline"] = process.Cmdline

		row := table.NewRow(rowData).WithStyle(standardRowStyle)
		rows = append(rows, row)
	}

	return rows
}
