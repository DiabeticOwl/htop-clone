package main

import (
	"os/exec"
	"strconv"
	"strings"
)

func getProcessesInfo() []processInfo {
	var processes []processInfo

	cmd := "ps"

	// Arguments for the command. Each comma adds a space to the output.
	// ax   : Lists all processes in the system.
	// o    : Gives a specific format to each process.
	// pid  : Process ID.
	// comm : Name of the process or Command.
	// pcpu : Percentage of the CPU used by the process.
	// prio : Priority assigned to the process.
	// exe  : Path to the executable.
	// args : The full command of the process with all it's arguments.

	// Each number preceded by a semicolon describes the total length of the
	// attribute extracted. Length is then used for slicing the desired values.
	args := []string{"-axo", "pid:10,user:50,comm:50,pcpu:4,pri:2,exe:100,args"}
	output, err := exec.Command(cmd, args...).Output()
	if err != nil {
		panic(err)
	}

	// The first line is the column names.
	// Removes the last newline and splits the entire string by the remaining.
	// The result is a process's info in the given format.
	procStrings := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")[1:]
	for _, line := range procStrings {
		pId, err := strconv.Atoi(strings.TrimSpace(line[:10]))
		if err != nil {
			panic(err)
		}
		cpuP, err := strconv.ParseFloat(strings.TrimSpace(line[113:117]), 32)
		if err != nil {
			panic(err)
		}
		prio, err := strconv.Atoi(strings.TrimSpace(line[118:120]))
		if err != nil {
			panic(err)
		}

		process := processInfo{
			PId:           int32(pId),
			User:          strings.TrimSpace(line[11:61]),
			Name:          strings.TrimSpace(line[62:112]),
			Priority:      int32(prio),
			CpuPercentage: cpuP,
			Cmdline:       strings.TrimSpace(line[222:]),
			ExeP:          strings.TrimSpace(line[121:221]),
		}

		processes = append(processes, process)
	}

	return processes
}
