package main

import (
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// TODO: Determine which constant better suits a given value.
const (
	_ = iota // Underscore means that it will be ignored.
	// This creates a block of constants that will return
	// 1 * 2^(10*iota). This formula is the description of
	// bytes in the binary system used by computers.
	KB = 1 << (10 * iota)
	MB
	GB
)

type processInfo struct {
	PId           int32
	User          string
	Name          string
	Priority      int32
	Niceness      int32
	CpuPercentage float64
	Cmdline       string
}

type diskInfo struct {
	Device    string
	MountPath string
	TotalSize float64
	FreeSize  float64
	UsedSize  float64
}

type virtualMemoryInfo struct {
	Total       float64
	Used        float64
	Available   float64
	UsedPercent float64
}

type swapMemoryInfo struct {
	Total       float64
	Used        float64
	Free        float64
	UsedPercent float64
}

func extractCpuInfo(interval time.Duration) []float64 {
	cpuInfo, _ := cpu.Percent(0, true)
	return cpuInfo
}

// extractMemoryInfo returns virtual and swap memory.
func extractMemoryInfo() (virtualMemoryInfo, swapMemoryInfo) {
	vm, _ := mem.VirtualMemory()
	sm, _ := mem.SwapMemory()

	vM := virtualMemoryInfo{
		Total:       float64(vm.Total) / GB,
		Used:        float64(vm.Used) / GB,
		Available:   float64(vm.Available) / GB,
		UsedPercent: vm.UsedPercent,
	}

	sM := swapMemoryInfo{
		Total:       float64(sm.Total) / GB,
		Used:        float64(sm.Used) / GB,
		Free:        float64(sm.Free) / GB,
		UsedPercent: sm.UsedPercent,
	}

	return vM, sM
}

func extractDiskInfo() []diskInfo {
	var disks []diskInfo

	dps, _ := disk.Partitions(true)
	for _, dsk := range dps {
		mount := dsk.Mountpoint
		if mount == "/" || strings.HasPrefix(mount, "/media/") {
			dskUsg, _ := disk.Usage(mount)

			disks = append(disks, diskInfo{
				Device:    dsk.Device,
				MountPath: mount,
				TotalSize: float64(dskUsg.Total) / GB,
				FreeSize:  float64(dskUsg.Free) / GB,
				UsedSize:  float64(dskUsg.Used) / GB,
			})
		}
	}

	return disks
}

func extractProcessesInfo() []processInfo {
	ps, _ := process.Processes()
	var processes []processInfo

	for _, p := range ps {
		u, _ := p.Username()
		n, _ := p.Name()
		prio, _ := p.IOnice()
		nice, _ := p.Nice()
		cPcg, _ := p.CPUPercent()
		cmdL, _ := p.Cmdline()

		if n == "spotify" {
			processes = append(processes, processInfo{
				PId:           p.Pid,
				User:          u,
				Name:          n,
				Priority:      prio,
				Niceness:      nice,
				CpuPercentage: cPcg,
				Cmdline:       cmdL,
			})
		}
	}

	return processes
}
