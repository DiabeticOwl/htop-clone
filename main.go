package main

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

const (
	_ = iota // Underscore means that it will be ignored.
	// This creates a block of constants that will return
	// 1 * 2^(10*iota). This formula is the description of
	// bytes in the binary system used by computers.
	KB = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func main() {
	fmt.Printf("\n------------ Memory info. ------------\n")
	vm, _ := mem.VirtualMemory()
	sm, _ := mem.SwapMemory()

	fmt.Printf("Total: %.2f GB, Used: %.2f GB, Available: %.2f GB, UsedPercent: %.4f%%\n", float64(vm.Total)/GB, float64(vm.Used)/GB, float64(vm.Available)/GB, vm.UsedPercent)
	fmt.Printf("Total: %.2f GB, Used: %.2f GB, Free: %.2f GB, UsedPercent: %.4f%%\n", float64(sm.Total)/GB, float64(sm.Used)/GB, float64(sm.Free)/GB, sm.UsedPercent)
	fmt.Printf("\n------------ Memory info. ------------\n")

	fmt.Printf("\n------------ Processes info. ------------\n")
	fmt.Printf("\n---- PID | Name | CPU Percentage ----\n")
	ps, _ := process.Processes()

	// TODO: Filter all processes as they normally are many.
	for _, p := range ps {
		n, _ := p.Name()
		pcg, _ := p.CPUPercent()

		fmt.Printf("\n%v | %s | %.4f%%\n", p.Pid, n, pcg)
	}
	fmt.Printf("\n------------ Processes info. ------------\n")

	fmt.Printf("\n------------ Disk info. ------------\n")
	dps, _ := disk.Partitions(true)
	for _, dsk := range dps {
		if mount := dsk.Mountpoint; mount == "/" || strings.HasPrefix(mount, "/media/") {
			dskUsg, _ := disk.Usage(mount)
			fmt.Printf("\n%v\t%s\t%.2f GB\t%.2f GB\t%.2f GB", dsk.Device, mount, float64(dskUsg.Total)/GB, float64(dskUsg.Free)/GB, float64(dskUsg.Used)/GB)
		}
	}
	fmt.Printf("\n------------ Disk info. ------------\n")
}
