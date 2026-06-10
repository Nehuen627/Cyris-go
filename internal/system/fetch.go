package system

import (
	"fmt"

	"cyris/internal/structs"

	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetHardwareInfo() structs.SystemSpecs {
	var err error

	var cpuInfo []cpu.InfoStat
	cpuInfo, err = cpu.Info()
	if err != nil {
		fmt.Println("CPU info error:", err.Error())
	}

	var cpuCores int
	cpuCores, err = cpu.Counts(true)
	if err != nil {
		fmt.Println("CPU cores error:", err.Error())
	}

	var ramInfo *mem.VirtualMemoryStat
	ramInfo, err = mem.VirtualMemory()
	if err != nil {
		fmt.Println("RAM error:", err.Error())
	}

	var diskInfo *disk.UsageStat
	diskInfo, err = disk.Usage("/")
	if err != nil {
		fmt.Println("Disk error:", err.Error())
	}

	var osInfo *host.InfoStat
	osInfo, err = host.Info()
	if err != nil {
		fmt.Println("OS info error:", err.Error())
	}

	var gpuName, gpuVendor string
	gpuInfo, err := ghw.GPU()
	if err != nil {
		fmt.Println("GPU error:", err.Error())
	} else if len(gpuInfo.GraphicsCards) > 0 {
		card := gpuInfo.GraphicsCards[0]
		gpuName = card.DeviceInfo.Product.Name
		gpuVendor = card.DeviceInfo.Vendor.Name
	}

	var cpuName string
	if len(cpuInfo) > 0 {
		cpuName = cpuInfo[0].ModelName
	}

	return structs.SystemSpecs{
		CPUName:   cpuName,
		CPUCores:  cpuCores,
		RAMTotal:  ramInfo.Total / 1024 / 1024,
		DiskFree:  diskInfo.Free / 1024 / 1024 / 1024,
		OS:        osInfo.OS,
		GPUName:   gpuName,
		GPUVendor: gpuVendor,
	}
}
