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

type NetworkStats struct {
	ActiveHalfOpenAttempts      int   `json:"activehalfopenattempts"`
	BytesWritten                int64 `json:"byteswritten"`
	BytesWrittenData            int64 `json:"byteswrittendata"`
	BytesRead                   int64 `json:"bytesread"`
	BytesReadData               int64 `json:"bytesreaddata"`
	BytesReadUsefulData         int64 `json:"bytesreadusefuldata"`
	BytesReadUsefulIntendedData int64 `json:"bytesreadusefulintendeddata"`
	ChunksWritten               int64 `json:"chunkswritten"`
	ChunksRead                  int64 `json:"chunksread"`
	ChunksReadUseful            int64 `json:"chunksreaduseful"`
	ChunksReadWasted            int64 `json:"chunksreadwasted"`
	MetadataChunksRead          int64 `json:"metadatachunksread"`
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
	if cpuprct, err := cpu.Percent(0, false); err == nil {
		if len(cpuprct) > 0 {
			s.CPU = cpuprct[0]
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

func GetNetworkStats() (retnstats NetworkStats) {
	s := Engine.Torc.Stats()
	retnstats.ActiveHalfOpenAttempts = s.ActiveHalfOpenAttempts
	retnstats.BytesWritten = s.BytesWritten.Int64()
	retnstats.BytesWrittenData = s.BytesWrittenData.Int64()
	retnstats.BytesRead = s.BytesRead.Int64()
	retnstats.BytesReadData = s.BytesReadData.Int64()
	retnstats.BytesReadUsefulData = s.BytesReadUsefulData.Int64()
	retnstats.BytesReadUsefulIntendedData = s.BytesReadUsefulIntendedData.Int64()
	retnstats.ChunksWritten = s.ChunksWritten.Int64()
	retnstats.ChunksRead = s.ChunksRead.Int64()
	retnstats.ChunksReadUseful = s.ChunksReadUseful.Int64()
	retnstats.ChunksReadWasted = s.ChunksReadWasted.Int64()
	retnstats.MetadataChunksRead = s.MetadataChunksRead.Int64()
	return
}
