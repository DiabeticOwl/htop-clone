package main

import (
	"fmt"
	"math"
	"runtime"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	cpuTableTitle           = "CPU Usage Percentage"
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
			Foreground(lipgloss.Color("#c1d0e8")).
			BorderBackground(lipgloss.Color("#7a89a3")).
			Align(lipgloss.Center))

	standardRowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
)

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

func generateCpuTableRows(m model) []table.Row {
	valStr := "CPU #%d: %s"
	rowCount := int(math.Ceil(float64(len(m.CpuInfo)) / cpuTableMaxColumnAmount))

	if rowCount*cpuTableMaxColumnAmount != len(m.CpuInfo) {
		s := "The amount of columns or rows is incorrect.\n"
		s = fmt.Sprintf("%sThe product of both should be %d.", s, len(m.CpuInfo))

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
				r += fmt.Sprintf("%-20s ", fmt.Sprintf(valStr, index, m.cpuProgresses[index].ViewAs(m.CpuInfo[i]/100)))
			} else {
				r += fmt.Sprintf("%-20s ", fmt.Sprintf(valStr, i, m.cpuProgresses[i].ViewAs(m.CpuInfo[i]/100)))
			}
		}

		nRow := table.NewRow(table.RowData{
			columnKeyCpuTable: r,
		}).WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(strconv.Itoa(255 - i*3))))
		rows = append(rows, nRow)
	}

	return rows
}

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

func generateMemoryTableRows(m model) []table.Row {
	vMemoryProg := m.memoryProgresses[0].ViewAs(m.VMemoryInfo.UsedPercent / 100)
	vMemoryView := fmt.Sprintf("%s %.2f GB/%.2f GB", vMemoryProg, m.VMemoryInfo.Used, m.VMemoryInfo.Total)

	sMemoryProg := m.memoryProgresses[1].ViewAs(m.SMemoryInfo.UsedPercent / 100)
	sMemoryView := fmt.Sprintf("%s %.2f GB/%.2f GB", sMemoryProg, m.SMemoryInfo.Used, m.SMemoryInfo.Total)

	rows := []table.Row{
		table.NewRow(table.RowData{
			columnKeyVirtualMemory: vMemoryView,
			columnKeySwapMemory:    sMemoryView,
		}),
	}

	return rows
}

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
		WithBaseStyle(styleBase).
		WithTargetWidth(m.Width).
		SortByAsc("FsType").
		ThenSortByAsc("MountPath")
}

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
