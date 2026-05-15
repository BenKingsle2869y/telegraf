package system_health

import (
	"fmt"
	"runtime"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemHealth collects basic system health metrics including CPU and memory usage.
type SystemHealth struct {
	CollectCPU    bool    `toml:"collect_cpu"`
	CollectMemory bool    `toml:"collect_memory"`
	CPUWarning    float64 `toml:"cpu_warning_threshold"`
	MemWarning    float64 `toml:"mem_warning_threshold"`

	Log telegraf.Logger `toml:"-"`
}

const sampleConfig = `
  ## Collect CPU usage metrics
  collect_cpu = true

  ## Collect memory usage metrics
  collect_memory = true

  ## CPU usage warning threshold (percentage)
  cpu_warning_threshold = 90.0

  ## Memory usage warning threshold (percentage)
  mem_warning_threshold = 90.0
`

func (s *SystemHealth) SampleConfig() string {
	return sampleConfig
}

func (s *SystemHealth) Description() string {
	return "Collects basic system health metrics (CPU and memory usage)"
}

func (s *SystemHealth) Gather(acc telegraf.Accumulator) error {
	now := time.Now()

	if s.CollectCPU {
		percents, err := cpu.Percent(0, false)
		if err != nil {
			return fmt.Errorf("error collecting CPU metrics: %w", err)
		}
		if len(percents) > 0 {
			tags := map[string]string{"host": runtime.GOOS}
			fields := map[string]interface{}{
				"usage_percent": percents[0],
				"warning":       percents[0] >= s.CPUWarning,
			}
			acc.AddGauge("system_health_cpu", fields, tags, now)
		}
	}

	if s.CollectMemory {
		vmStat, err := mem.VirtualMemory()
		if err != nil {
			return fmt.Errorf("error collecting memory metrics: %w", err)
		}
		tags := map[string]string{"host": runtime.GOOS}
		fields := map[string]interface{}{
			"used_percent": vmStat.UsedPercent,
			"total_bytes":  vmStat.Total,
			"used_bytes":   vmStat.Used,
			"free_bytes":   vmStat.Free,
			"warning":      vmStat.UsedPercent >= s.MemWarning,
		}
		acc.AddGauge("system_health_memory", fields, tags, now)
	}

	return nil
}

func init() {
	inputs.Add("system_health", func() telegraf.Input {
		return &SystemHealth{
			CollectCPU:    true,
			CollectMemory: true,
			CPUWarning:    90.0,
			MemWarning:    90.0,
		}
	})
}
