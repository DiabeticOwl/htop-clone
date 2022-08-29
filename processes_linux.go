package main

import (
	"os/exec"
	"strconv"
	"strings"
)

func extractProcessesInfo() []processInfo {
	var processes []processInfo

	cmd := "ps"
	args := []string{"-axo", "pid:10,user:50,comm:50,pcpu:4,pri:2,exe:100,args"}
	output, err := exec.Command(cmd, args...).Output()
	if err != nil {
		panic(err)
	}

	procStrings := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")[1:]
	for _, line := range procStrings {
		pId, err := (strconv.ParseInt(strings.TrimSpace(line[:10]), 10, 32))
		if err != nil {
			panic(err)
		}
		cpuP, err := strconv.ParseFloat(strings.TrimSpace(line[113:117]), 32)
		if err != nil {
			panic(err)
		}
		prio, err := (strconv.ParseInt(strings.TrimSpace(line[118:120]), 10, 32))
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
