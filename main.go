package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v3/cpu"
)

// TODO: Determine which constant is better suits a given value.
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
	ExePath       string
}

type diskInfo struct {
	Device    string
	MountPath string
	TotalSize uint64
	FreeSize  uint64
	UsedSize  uint64
}

func extractCpuInfo(interval time.Duration) []float64 {
	cpuInfo, _ := cpu.Percent(0, true)
	return cpuInfo
}

func main() {
	p := tea.NewProgram(model{})
	if err := p.Start(); err != nil {
		panic(err)
	}
	// fmt.Printf("\n------------ Memory info. ------------\n")
	// vm, _ := mem.VirtualMemory()
	// sm, _ := mem.SwapMemory()

	// fmt.Printf("Total: %.2f GB, Used: %.2f GB, Available: %.2f GB, UsedPercent: %.4f%%\n", float64(vm.Total)/GB, float64(vm.Used)/GB, float64(vm.Available)/GB, vm.UsedPercent)
	// fmt.Printf("Total: %.2f GB, Used: %.2f GB, Free: %.2f GB, UsedPercent: %.4f%%\n", float64(sm.Total)/GB, float64(sm.Used)/GB, float64(sm.Free)/GB, sm.UsedPercent)
	// fmt.Printf("\n------------ Memory info. ------------\n")

	// fmt.Printf("\n------------ Processes info. ------------\n")
	// fmt.Printf("\n---- PID | Name | CPU Percentage ----\n")
	// ps, _ := process.Processes()
	// var processes []processInfo

	// for _, p := range ps {
	// 	u, _ := p.Username()
	// 	n, _ := p.Name()
	// 	prio, _ := p.IOnice()
	// 	nice, _ := p.Nice()
	// 	cPcg, _ := p.CPUPercent()
	// 	exeP, _ := p.Exe()

	// 	processes = append(processes, processInfo{
	// 		PId:           p.Pid,
	// 		User:          u,
	// 		Name:          n,
	// 		Priority:      prio,
	// 		Niceness:      nice,
	// 		CpuPercentage: cPcg,
	// 		ExePath:       exeP,
	// 	})

	// 	// if n == "spotify" {
	// 	// 	fmt.Printf("\n%v | %s | %.4f%%\n", p.Pid, n, pcg)
	// 	// }
	// }
	// fmt.Printf("\n------------ Processes info. ------------\n")

	// fmt.Printf("\n------------ Disk info. ------------\n")
	// var disks []diskInfo
	// dps, _ := disk.Partitions(true)
	// for _, dsk := range dps {
	// 	if mount := dsk.Mountpoint; mount == "/" || strings.HasPrefix(mount, "/media/") {
	// 		dskUsg, _ := disk.Usage(mount)

	// 		disks = append(disks, diskInfo{
	// 			Device:    dsk.Device,
	// 			MountPath: mount,
	// 			TotalSize: dskUsg.Total,
	// 			FreeSize:  dskUsg.Free,
	// 			UsedSize:  dskUsg.Used,
	// 		})

	// 		fmt.Printf("\n%v\t%s\t%.2f GB\t%.2f GB\t%.2f GB", dsk.Device, mount, float64(dskUsg.Total)/GB, float64(dskUsg.Free)/GB, float64(dskUsg.Used)/GB)
	// 	}
	// }
	// fmt.Printf("\n------------ Disk info. ------------\n")
}
