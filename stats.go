package main

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
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
		"apfs":    {},
	}
)

type memoryInfo struct {
	Total       float64
	Used        float64
	UsedPercent float64
}

type diskInfo struct {
	FsType    string
	Device    string
	MountPath string
	TotalSize float64
	FreeSize  float64
	UsedSize  float64
}

type processInfo struct {
	PId           int32
	User          string
	Name          string
	Priority      int32
	CpuPercentage float64
	Cmdline       string
	ExeP          string
}

func extractCpuInfo() []float64 {
	cpuInfo, _ := cpu.Percent(0, true)
	return cpuInfo
}

// extractMemoryInfo returns virtual and swap memory.
func extractMemoryInfo() (memoryInfo, memoryInfo) {
	vm, _ := mem.VirtualMemory()
	sm, _ := mem.SwapMemory()

	vMemoryInfo := memoryInfo{
		Total:       float64(vm.Total) / GB,
		Used:        float64(vm.Used) / GB,
		UsedPercent: vm.UsedPercent,
	}

	sMemoryInfo := memoryInfo{
		Total:       float64(sm.Total) / GB,
		Used:        float64(sm.Used) / GB,
		UsedPercent: sm.UsedPercent,
	}

	return vMemoryInfo, sMemoryInfo
}

func extractDiskInfo() []diskInfo {
	var disks []diskInfo

	dps, _ := disk.Partitions(true)
	for _, dsk := range dps {
		if _, ok := desiredFileSystems[dsk.Fstype]; ok {
			mount := dsk.Mountpoint
			dskUsg, _ := disk.Usage(mount)

			diskInfo := diskInfo{
				FsType:    dsk.Fstype,
				Device:    dsk.Device,
				MountPath: mount,
				TotalSize: float64(dskUsg.Total) / GB,
				FreeSize:  float64(dskUsg.Free) / GB,
				UsedSize:  float64(dskUsg.Used) / GB,
			}

			disks = append(disks, diskInfo)
		}
	}

	return disks
}
