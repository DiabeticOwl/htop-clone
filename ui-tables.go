package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	tableDefaultTargetWidth = 70
	tableLargeTargetWidth   = 160
	tableLargestTargetWidth = 250
	// https://pkg.go.dev/github.com/evertras/bubble-table@v0.14.4/table?utm_source=gopls#NewFlexColumn
	columnDefaultFlexFactor = 1

	// * CPU Table *

	cpuTableTitle           = "CPU Usage Percentage"
	cpuTableMaxColumnAmount = 3

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
			Foreground(lipgloss.Color("#a7a")).
			BorderBackground(lipgloss.Color("#a38")).
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
		"Niceness":      "Nice",
		"CpuPercentage": "CPU Usage Percentage",
		"Name":          "Name",
		"ExeP":          "Executable Path",
		"Cmdline":       "Command",
	}
	// * Processes Table *
)

// * CPU Table *

func newCpuTable() table.Model {
	columns := []table.Column{
		table.NewFlexColumn(columnKeyCpuTable, cpuTableTitle,
			columnDefaultFlexFactor),
	}

	return (table.
		New(columns).
		BorderRounded().
		WithBaseStyle(styleBase).
		WithTargetWidth(tableLargestTargetWidth))
}

func generateCpuTableRows(m model) []table.Row {
	valStr := "CPU #%d: %s"
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
				r += fmt.Sprintf("%-20s ", fmt.Sprintf(valStr, index, m.progresses[index].ViewAs(m.CpuInfo[i]/100)))
			} else {
				r += fmt.Sprintf("%-20s ", fmt.Sprintf(valStr, i, m.progresses[i].ViewAs(m.CpuInfo[i]/100)))
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

func newMemoryTable() table.Model {
	columns := []table.Column{
		table.NewFlexColumn(columnKeyVirtualMemory, columnKeyVirtualMemoryTitle,
			columnDefaultFlexFactor),
		table.NewFlexColumn(columnKeySwapMemory, columnKeySwapMemoryTitle,
			columnDefaultFlexFactor),
	}

	return (table.
		New(columns).
		BorderRounded().
		WithBaseStyle(styleBase).
		WithTargetWidth(tableDefaultTargetWidth))
}

func generateMemoryTableRows(m model) []table.Row {
	gbFormat := "%s: %.2f GB"
	pcFormat := "%s: %.4f%%"

	totalM := table.NewRow(table.RowData{
		columnKeyVirtualMemory: fmt.Sprintf(gbFormat, "Total", m.VMemoryInfo["Total"]),
		columnKeySwapMemory:    fmt.Sprintf(gbFormat, "Total", m.SMemoryInfo["Total"]),
	}).WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("255")))
	usedM := table.NewRow(table.RowData{
		columnKeyVirtualMemory: fmt.Sprintf(gbFormat, "Used", m.VMemoryInfo["Used"]),
		columnKeySwapMemory:    fmt.Sprintf(gbFormat, "Used", m.SMemoryInfo["Used"]),
	}).WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("252")))
	availableM := table.NewRow(table.RowData{
		columnKeyVirtualMemory: fmt.Sprintf(gbFormat, "Available", m.VMemoryInfo["Available"]),
		columnKeySwapMemory:    fmt.Sprintf(gbFormat, "Free", m.SMemoryInfo["Free"]),
	}).WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("249")))
	usedPcM := table.NewRow(table.RowData{
		columnKeyVirtualMemory: fmt.Sprintf(pcFormat, "UsedPercent", m.VMemoryInfo["UsedPercent"]),
		columnKeySwapMemory:    fmt.Sprintf(pcFormat, "UsedPercent", m.SMemoryInfo["UsedPercent"]),
	}).WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("246")))

	rows := []table.Row{totalM, usedM, availableM, usedPcM}

	return rows
}

// * Memory Table *

// * Disks Table *

func newDisksTable() table.Model {
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

	return (table.
		New(columns).
		BorderRounded().
		WithBaseStyle(styleBase).
		WithTargetWidth(tableLargeTargetWidth).
		SortByAsc("FsType").
		ThenSortByAsc("MountPath"))
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

func newProcessesTable(pCount int) table.Model {
	var columns []table.Column
	columnsOrder := []string{"PId", "User", "Priority", "Niceness",
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

	return (table.
		New(columns).
		BorderRounded().
		WithBaseStyle(styleBase).
		WithTargetWidth(tableLargestTargetWidth).
		WithPageSize(pCount).
		SortByDesc("CpuPercentage").
		Focused(true))
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
