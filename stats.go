package main

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	_ = iota
	// Bit shift to get constants of 1,024 bytes, 1,048,576 bytes, etc.
	// The formula is: 1 * 2^(10*iota).
	KB = 1 << (10 * iota)
	MB
	GB
)

var (
	// Displaying relevant file systems. Irrelevant may be
	// read only (like squashfs), tracing (like tracefs), etc.
	fsFilter = map[string]struct{}{
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

// getCpuInfo returns the information of the cores in the system.
func getCpuInfo() []float64 {
	cpuInfo, _ := cpu.Percent(0, true)
	return cpuInfo
}

// getMemoryInfo returns virtual and swap memory.
func getMemoryInfo() (memoryInfo, memoryInfo) {
	// Ignoring errors because of unaccounted ones in these methods.
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

// getDiskInfo returns an array of the information of relevant disks in the
// system.
func getDiskInfo() []diskInfo {
	var disks []diskInfo

	dps, _ := disk.Partitions(true)
	for _, dsk := range dps {
		if _, ok := fsFilter[dsk.Fstype]; !ok {
			continue
		}

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

	return disks
}
