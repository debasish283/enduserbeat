// +build darwin linux windows
// +build cgo

package instance

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/satori/go.uuid"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/metric/system/cpu"
	"github.com/elastic/beats/libbeat/metric/system/process"
	"github.com/elastic/beats/libbeat/monitoring"
	"github.com/elastic/beats/libbeat/monitoring/report/log"
)

var (
	beatProcessStats *process.Stats
	ephemeralID      uuid.UUID
)

func init() {
	beatMetrics := monitoring.Default.NewRegistry("beat")
	monitoring.NewFunc(beatMetrics, "memstats", reportMemStats, monitoring.Report)
	monitoring.NewFunc(beatMetrics, "cpu", reportBeatCPU, monitoring.Report)
	monitoring.NewFunc(beatMetrics, "info", reportInfo, monitoring.Report)

	systemMetrics := monitoring.Default.NewRegistry("system")
	monitoring.NewFunc(systemMetrics, "cpu", reportSystemCPUUsage, monitoring.Report)
	if runtime.GOOS != "windows" {
		monitoring.NewFunc(systemMetrics, "load", reportSystemLoadAverage, monitoring.Report)
	}

	ephemeralID = uuid.NewV4()
}

func setupMetrics(name string) error {
	beatProcessStats = &process.Stats{
		Procs:        []string{name},
		EnvWhitelist: nil,
		CpuTicks:     true,
		CacheCmdLine: true,
		IncludeTop:   process.IncludeTopConfig{},
	}
	err := beatProcessStats.Init()

	return err
}

func reportMemStats(m monitoring.Mode, V monitoring.Visitor) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	V.OnRegistryStart()
	defer V.OnRegistryFinished()

	monitoring.ReportInt(V, "memory_total", int64(stats.TotalAlloc))
	if m == monitoring.Full {
		monitoring.ReportInt(V, "memory_alloc", int64(stats.Alloc))
		monitoring.ReportInt(V, "gc_next", int64(stats.NextGC))
	}

	rss, err := getRSSSize()
	if err != nil {
		logp.Err("Error while getting memory usage: %v", err)
		return
	}
	monitoring.ReportInt(V, "rss", int64(rss))
}

func getRSSSize() (uint64, error) {
	beatPID := os.Getpid()
	state, err := beatProcessStats.GetOne(beatPID)
	if err != nil {
		return 0, fmt.Errorf("error retrieving process stats")
	}

	iRss, err := state.GetValue("memory.rss.bytes")
	if err != nil {
		return 0, fmt.Errorf("error getting Resident Set Size: %v", err)
	}

	rss, ok := iRss.(uint64)
	if !ok {
		return 0, fmt.Errorf("error converting Resident Set Size: %v", err)
	}
	return rss, nil
}

func reportInfo(_ monitoring.Mode, V monitoring.Visitor) {
	V.OnRegistryStart()
	defer V.OnRegistryFinished()

	delta := time.Since(log.StartTime)
	uptime := int64(delta / time.Millisecond)
	monitoring.ReportNamespace(V, "uptime", func() {
		monitoring.ReportInt(V, "ms", uptime)
	})

	monitoring.ReportString(V, "ephemeral_id", ephemeralID.String())
}

func reportBeatCPU(_ monitoring.Mode, V monitoring.Visitor) {
	V.OnRegistryStart()
	defer V.OnRegistryFinished()

	totalCPUUsage, cpuTicks, err := getCPUUsage()
	if err != nil {
		logp.Err("Error retrieving CPU percentages: %v", err)
		return
	}

	userTime, systemTime, err := process.GetOwnResourceUsageTimeInMillis()
	if err != nil {
		logp.Err("Error retrieving CPU usage time: %v", err)
		return
	}

	monitoring.ReportNamespace(V, "user", func() {
		monitoring.ReportInt(V, "ticks", int64(cpuTicks.User))
		monitoring.ReportInt(V, "time", userTime)
	})
	monitoring.ReportNamespace(V, "system", func() {
		monitoring.ReportInt(V, "ticks", int64(cpuTicks.System))
		monitoring.ReportInt(V, "time", systemTime)
	})
	monitoring.ReportNamespace(V, "total", func() {
		monitoring.ReportFloat(V, "value", totalCPUUsage)
		monitoring.ReportInt(V, "ticks", int64(cpuTicks.Total))
		monitoring.ReportInt(V, "time", userTime+systemTime)
	})
}

func getCPUUsage() (float64, *process.Ticks, error) {
	beatPID := os.Getpid()
	state, err := beatProcessStats.GetOne(beatPID)
	if err != nil {
		return 0.0, nil, fmt.Errorf("error retrieving process stats")
	}

	iTotalCPUUsage, err := state.GetValue("cpu.total.value")
	if err != nil {
		return 0.0, nil, fmt.Errorf("error getting total CPU since start: %v", err)
	}

	totalCPUUsage, ok := iTotalCPUUsage.(float64)
	if !ok {
		return 0.0, nil, fmt.Errorf("error converting value of CPU usage since start")
	}

	iTotalCPUUserTicks, err := state.GetValue("cpu.user.ticks")
	if err != nil {
		return 0.0, nil, fmt.Errorf("error getting number of user CPU ticks since start: %v", err)
	}

	totalCPUUserTicks, ok := iTotalCPUUserTicks.(uint64)
	if !ok {
		return 0.0, nil, fmt.Errorf("error converting value of user CPU ticks since start")
	}

	iTotalCPUSystemTicks, err := state.GetValue("cpu.system.ticks")
	if err != nil {
		return 0.0, nil, fmt.Errorf("error getting number of system CPU ticks since start: %v", err)
	}

	totalCPUSystemTicks, ok := iTotalCPUSystemTicks.(uint64)
	if !ok {
		return 0.0, nil, fmt.Errorf("error converting value of system CPU ticks since start")
	}

	iTotalCPUTicks, err := state.GetValue("cpu.total.ticks")
	if err != nil {
		return 0.0, nil, fmt.Errorf("error getting total number of CPU ticks since start: %v", err)
	}

	totalCPUTicks, ok := iTotalCPUTicks.(uint64)
	if !ok {
		return 0.0, nil, fmt.Errorf("error converting total value of CPU ticks since start")
	}

	p := process.Ticks{
		User:   totalCPUUserTicks,
		System: totalCPUSystemTicks,
		Total:  totalCPUTicks,
	}

	return totalCPUUsage, &p, nil
}

func reportSystemLoadAverage(_ monitoring.Mode, V monitoring.Visitor) {
	V.OnRegistryStart()
	defer V.OnRegistryFinished()

	load, err := cpu.Load()
	if err != nil {
		logp.Err("Error retrieving load average: %v", err)
		return
	}
	avgs := load.Averages()
	monitoring.ReportFloat(V, "1", avgs.OneMinute)
	monitoring.ReportFloat(V, "5", avgs.FiveMinute)
	monitoring.ReportFloat(V, "15", avgs.FifteenMinute)

	normAvgs := load.NormalizedAverages()
	monitoring.ReportNamespace(V, "norm", func() {
		monitoring.ReportFloat(V, "1", normAvgs.OneMinute)
		monitoring.ReportFloat(V, "5", normAvgs.FiveMinute)
		monitoring.ReportFloat(V, "15", normAvgs.FifteenMinute)
	})
}

func reportSystemCPUUsage(_ monitoring.Mode, V monitoring.Visitor) {
	V.OnRegistryStart()
	defer V.OnRegistryFinished()

	monitoring.ReportInt(V, "cores", int64(process.NumCPU))
}
