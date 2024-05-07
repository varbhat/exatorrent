package core

import (
	"os"
	"runtime"
	"time"

	"github.com/pbnjay/memory"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type machInfo struct {
	Arch       string    `json:"arch"`
	NumberCPUs int       `json:"numbercpu"`
	CPUModel   string    `json:"cpumodel"`
	HostName   string    `json:"hostname"`
	Platform   string    `json:"platform"`
	OS         string    `json:"os"`
	TotalMem   uint64    `json:"totalmem"`
	GoVersion  string    `json:"goversion"`
	StartedAt  time.Time `json:"startedat"`
}

type machStats struct {
	CPU             float64 `json:"cpucycles"`
	DiskFree        uint64  `json:"diskfree"`
	DiskUsedPercent float64 `json:"diskpercent"`
	MemUsedPercent  float64 `json:"mempercent"`
	GoMemory        int64   `json:"gomem"`
	GoMemorySys     int64   `json:"gomemsys"`
	GoRoutines      int     `json:"goroutines"`
}

var MachInfo machInfo = loadMachInfo()
var MachStats machStats

func loadMachInfo() (retmachinfo machInfo) {
	hostInfo, err := host.Info()
	if err == nil && hostInfo != nil {
		retmachinfo.OS = hostInfo.Platform + " " + hostInfo.PlatformVersion
		retmachinfo.Platform = hostInfo.OS
	}

	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		retmachinfo.CPUModel = cpuInfo[0].ModelName
	}

	retmachinfo.HostName, _ = os.Hostname()
	retmachinfo.GoVersion = runtime.Version()
	retmachinfo.TotalMem = memory.TotalMemory()
	retmachinfo.Arch = runtime.GOARCH
	retmachinfo.NumberCPUs = runtime.NumCPU()
	retmachinfo.StartedAt = time.Now()
	return retmachinfo
}

func (s *machStats) LoadStats(diskDir string) {
	//count cpu cycles between last count
	if cpu, err := cpu.Percent(0, false); err == nil {
		if len(cpu) > 0 {
			s.CPU = cpu[0]
		}
	}
	//count disk usage
	if stat, err := disk.Usage(diskDir); err == nil {
		s.DiskUsedPercent = stat.UsedPercent
		s.DiskFree = stat.Free
	}
	//count memory usage
	if stat, err := mem.VirtualMemory(); err == nil {
		s.MemUsedPercent = stat.UsedPercent
	}
	//count total bytes allocated by the go runtime
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)
	s.GoMemory = int64(memStats.Alloc)
	s.GoMemorySys = int64(memStats.Sys)
	//count current number of goroutines
	s.GoRoutines = runtime.NumGoroutine()
}
