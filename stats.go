package main

import (
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

var (
	desiredFileSystems = map[string]struct{}{
		"ext4":    {},
		"vfat":    {},
		"fuseblk": {},
		"ntfs":    {},
		"fat32":   {},
	}

	virtualMemoryInfo = make(map[string]interface{})
	swapMemoryInfo    = make(map[string]interface{})
)

func extractCpuInfo() []float64 {
	cpuInfo, _ := cpu.Percent(0, true)
	return cpuInfo
}

// extractMemoryInfo returns virtual and swap memory.
func extractMemoryInfo() (map[string]interface{}, map[string]interface{}) {
	vm, _ := mem.VirtualMemory()
	sm, _ := mem.SwapMemory()

	virtualMemoryInfo["Total"] = float64(vm.Total) / GB
	virtualMemoryInfo["Used"] = float64(vm.Used) / GB
	virtualMemoryInfo["Available"] = float64(vm.Available) / GB
	virtualMemoryInfo["UsedPercent"] = vm.UsedPercent

	swapMemoryInfo["Total"] = float64(sm.Total) / GB
	swapMemoryInfo["Used"] = float64(sm.Used) / GB
	swapMemoryInfo["Free"] = float64(sm.Free) / GB
	swapMemoryInfo["UsedPercent"] = sm.UsedPercent

	return virtualMemoryInfo, swapMemoryInfo
}

func extractDiskInfo() []map[string]interface{} {
	var disks []map[string]interface{}

	dps, _ := disk.Partitions(true)
	for _, dsk := range dps {
		if _, ok := desiredFileSystems[dsk.Fstype]; ok {
			mount := dsk.Mountpoint
			dskUsg, _ := disk.Usage(mount)

			diskInfo := make(map[string]interface{})

			diskInfo["FsType"] = dsk.Fstype
			diskInfo["Device"] = dsk.Device
			diskInfo["MountPath"] = mount
			diskInfo["TotalSize"] = float64(dskUsg.Total) / GB
			diskInfo["FreeSize"] = float64(dskUsg.Free) / GB
			diskInfo["UsedSize"] = float64(dskUsg.Used) / GB

			disks = append(disks, diskInfo)
		}
	}

	return disks
}

func extractProcessesInfo() []map[string]interface{} {
	ps, _ := process.Processes()
	var processes []map[string]interface{}

	for _, p := range ps {
		processInfo := make(map[string]interface{})

		u, _ := p.Username()
		n, _ := p.Name()
		prio, _ := p.Nice()
		nice, _ := p.IOnice()
		cPcg, _ := p.CPUPercent()
		exeP, _ := p.Exe()
		cmdL, _ := p.Cmdline()

		processInfo["PId"] = p.Pid
		processInfo["User"] = u
		processInfo["Name"] = n
		processInfo["Priority"] = prio
		processInfo["Niceness"] = nice
		processInfo["CpuPercentage"] = cPcg
		processInfo["Cmdline"] = cmdL
		processInfo["ExeP"] = exeP

		processes = append(processes, processInfo)
	}

	return processes
}
