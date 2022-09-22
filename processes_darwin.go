package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const (
	// Columns width.
	smallW = "1234567890"
	largeW = smallW + smallW + smallW + smallW + smallW
	hugeW  = largeW + largeW
)

func getProcessesInfo() []processInfo {
	var processes []processInfo

	cmd := "ps"

	// Arguments for ps. Each comma adds a space to the output.
	// ax      : Lists all processes in the system.
	// c       : Makes comm and command keywords the same. Used for slicing and cleaning purposes.
	// o       : Gives a specific format to each process.
	// pid     : Process ID.
	// comm    : Name of the process.
	// command : Name of the process. (Kept to maintain uniformity with other operating systems.)
	// pcpu    : Percentage of the CPU used by the process.
	// prio    : Priority assigned to the process.
	// args    : The full command of the process with all it's arguments.

	keywords := fmt.Sprintf("pid=%s,user=%s,comm=%s,pcpu,pri,command=%s,args", smallW, largeW, hugeW, hugeW)
	args := []string{"-axcro", keywords}

	output, err := exec.Command(cmd, args...).Output()
	if err != nil {
		panic(err)
	}

	// The first line is the column names
	procStrings := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")[1:]
	for _, line := range procStrings {
		pId, err := strconv.ParseInt(strings.TrimSpace(line[:10]), 10, 32)
		if err != nil {
			panic(err)
		}
		cpuP, err := strconv.ParseFloat(strings.TrimSpace(line[163:168]), 32)
		if err != nil {
			panic(err)
		}
		prio, err := strconv.ParseInt(strings.TrimSpace(line[169:172]), 10, 32)
		if err != nil {
			panic(err)
		}

		process := processInfo{
			PId:           int32(pId),
			User:          strings.TrimSpace(line[11:61]),
			Name:          strings.TrimSpace(line[62:162]),
			Priority:      int32(prio),
			CpuPercentage: cpuP,
			Cmdline:       strings.TrimSpace(line[274:]),
			ExeP:          strings.TrimSpace(line[173:273]),
		}

		processes = append(processes, process)
	}

	return processes
}
