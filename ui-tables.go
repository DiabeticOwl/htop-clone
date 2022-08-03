package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	// https://pkg.go.dev/github.com/evertras/bubble-table@v0.14.4/table?utm_source=gopls#NewFlexColumn
	columnDefaultFlexFactor = 1

	// * CPU Table *

	cpuTableTitle           = "CPU Usage Percentage"
	cpuTableMaxColumnAmount = 4

	columnKeyCpuTable = "cpuTable"
	// * CPU Table *

	// * Memory Table *

	columnKeyVirtualMemory      = "virtualMemory"
	columnKeyVirtualMemoryTitle = "Virtual Memory"
	columnKeySwapMemory         = "swapMemory"
	columnKeySwapMemoryTitle    = "Swap Memory"
	// * Memory Table *
)

var (
	styleBase = (lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("#c1d0e8")).
			BorderBackground(lipgloss.Color("#7a89a3")).
			Align(lipgloss.Center))

	// * Disks Table *

	// ColumnKey: ColumnTitle
	disksColumnKeyMap = map[string]string{
		"FsType":    "File System Type",
		"Device":    "Device",
		"MountPath": "Mount Path",
		"TotalSize": "Total Size",
		"FreeSize":  "Free Size",
		"UsedSize":  "Used Size",
	}
	// * Disks Table *

	// * Processes Table *

	// ColumnKey: ColumnTitle
	processesColumnKeyMap = map[string]string{
		"PId":           "Process ID",
		"User":          "Username",
		"Priority":      "Priority",
		"CpuPercentage": "CPU Usage Percentage",
		"Name":          "Name",
		"ExeP":          "Executable Path",
		"Cmdline":       "Command",
	}
	// * Processes Table *
)

// * CPU Table *

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

// * CPU Table *

// * Memory Table *

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
	vMemoryProg := m.memoryProgresses[0].ViewAs(m.VMemoryInfo["UsedPercent"].(float64) / 100)
	vMemoryView := fmt.Sprintf("%s %.2f GB/%.2f GB", vMemoryProg, m.VMemoryInfo["Used"], m.VMemoryInfo["Total"])

	sMemoryProg := m.memoryProgresses[1].ViewAs(m.SMemoryInfo["UsedPercent"].(float64) / 100)
	sMemoryView := fmt.Sprintf("%s %.2f GB/%.2f GB", sMemoryProg, m.SMemoryInfo["Used"], m.SMemoryInfo["Total"])

	rows := []table.Row{
		table.NewRow(table.RowData{
			columnKeyVirtualMemory: vMemoryView,
			columnKeySwapMemory:    sMemoryView,
		}),
	}

	return rows
}

// * Memory Table *

// * Disks Table *

func newDisksTable(m model) table.Model {
	var columns []table.Column
	columnsOrder := []string{"FsType", "Device", "MountPath", "TotalSize",
		"FreeSize", "UsedSize"}

	for _, column := range columnsOrder {
		factor := columnDefaultFlexFactor
		if column == "MountPath" {
			factor *= 4
		}

		nCol := table.NewFlexColumn(column, disksColumnKeyMap[column], factor)

		columns = append(columns, nCol)
	}

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
	var format string
	rowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	for _, disk := range m.DisksInfo {
		rowData := make(table.RowData)
		for key, value := range disk {
			switch value.(type) {
			case float64:
				format = "%2.f GB"
			default:
				format = "%s"
			}

			rowData[key] = fmt.Sprintf(format, value)
		}

		row := table.NewRow(rowData).WithStyle(rowStyle)
		rows = append(rows, row)
	}

	return rows
}

// * Disks Table *

// * Processes Table *

func newProcessesTable(m model, pCount int) table.Model {
	var columns []table.Column
	columnsOrder := []string{"PId", "User", "Priority",
		"CpuPercentage", "Name", "ExeP", "Cmdline"}

	for _, column := range columnsOrder {
		factor := columnDefaultFlexFactor
		if column == "Name" || column == "CpuPercentage" || column == "User" {
			factor *= 2
		} else if column == "ExeP" {
			factor *= 4
		} else if column == "Cmdline" {
			factor *= 6
		}

		nCol := table.NewFlexColumn(column, processesColumnKeyMap[column], factor)

		columns = append(columns, nCol)
	}

	return table.
		New(columns).
		BorderRounded().
		WithBaseStyle(styleBase).
		WithTargetWidth(m.Width).
		WithPageSize(pCount).
		SortByDesc("CpuPercentage").
		Focused(true)
}

func generateProcessesTableRows(m model) []table.Row {
	var rows []table.Row
	var format string
	rowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	for _, process := range m.Processes {
		rowData := make(table.RowData)
		for key, value := range process {
			switch value.(type) {
			case int, int32:
				format = "%d"
			case float64:
				format = "%.4f%%"
			default:
				format = "%s"
			}

			rowData[key] = fmt.Sprintf(format, value)
		}

		row := table.NewRow(rowData).WithStyle(rowStyle)
		rows = append(rows, row)
	}

	return rows
}

// * Processes Table *
