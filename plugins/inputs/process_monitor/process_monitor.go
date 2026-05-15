package process_monitor

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

const sampleConfig = `
  ## List of process names to monitor
  processes = ["telegraf", "nginx", "postgres"]

  ## Collect metrics interval
  # interval = "10s"
`

type ProcessMonitor struct {
	Processes []string `toml:"processes"`
	Log       telegraf.Logger
}

func (pm *ProcessMonitor) SampleConfig() string {
	return sampleConfig
}

func (pm *ProcessMonitor) Description() string {
	return "Monitor CPU and memory usage of specific processes by name"
}

func (pm *ProcessMonitor) Gather(acc telegraf.Accumulator) error {
	for _, name := range pm.Processes {
		if err := pm.gatherProcess(acc, name); err != nil {
			pm.Log.Warnf("Error gathering process %s: %v", name, err)
		}
	}
	return nil
}

func (pm *ProcessMonitor) gatherProcess(acc telegraf.Accumulator, name string) error {
	out, err := exec.Command("pgrep", "-x", name).Output()
	if err != nil {
		// Process not found - report it as not running but don't treat as an error
		acc.AddFields("process_monitor",
			map[string]interface{}{"running": 0},
			map[string]string{"process": name},
			time.Now(),
		)
		return nil
	}

	pids := strings.Fields(strings.TrimSpace(string(out)))
	count := len(pids)

	var totalCPU, totalMem float64
	for _, pid := range pids {
		cpu, mem, err := getProcessStats(pid)
		if err != nil {
			continue
		}
		totalCPU += cpu
		totalMem += mem
	}

	// Calculate per-process averages when multiple PIDs are found
	var avgCPU, avgMem float64
	if count > 0 {
		avgCPU = totalCPU / float64(count)
		avgMem = totalMem / float64(count)
	}

	acc.AddFields("process_monitor",
		map[string]interface{}{
			"running":   1,
			"pid_count": count,
			"cpu_usage": totalCPU,
			"mem_usage": totalMem,
			"avg_cpu":   avgCPU,
			"avg_mem":   avgMem,
		},
		map[string]string{"process": name},
		time.Now(),
	)
	return nil
}

func getProcessStats(pid string) (float64, float64, error) {
	out, err := exec.Command("ps", "-p", pid, "-o", "%cpu,%mem", "--no-headers").Output()
	if err != nil {
		return 0, 0, fmt.Errorf("ps command failed: %w", err)
	}
	fields := strings.Fields(strings.TrimSpace(string(out)))
	if len(fields) < 2 {
		return 0, 0, fmt.Errorf("unexpected ps output")
	}
	cpu, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, 0, err
	}
	mem, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return 0, 0, err
	}
	return cpu, mem, nil
}

func init() {
	inputs.Add("process_monitor", func() telegraf.Input {
		return &ProcessMonitor{}
	})
}
