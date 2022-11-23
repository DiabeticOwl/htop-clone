package main

import (
	"github.com/shirou/gopsutil/v3/process"
)

func extractProcessesInfo() []processInfo {
	ps, _ := process.Processes()
	var processes []processInfo

	for _, p := range ps {
		u, _ := p.Username()
		n, _ := p.Name()
		prio, _ := p.Nice()
		cPcg, _ := p.CPUPercent()
		exeP, _ := p.Exe()
		cmdL, _ := p.Cmdline()

		processInfo := processInfo{
			PId:           p.Pid,
			User:          u,
			Name:          n,
			Priority:      prio,
			CpuPercentage: cPcg,
			Cmdline:       cmdL,
			ExeP:          exeP,
		}

		processes = append(processes, processInfo)
	}

	return processes
}
