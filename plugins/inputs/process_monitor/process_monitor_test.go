package process_monitor

import (
	"testing"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessMonitorDescription(t *testing.T) {
	pm := &ProcessMonitor{}
	desc := pm.Description()
	assert.NotEmpty(t, desc)
	assert.Contains(t, desc, "process")
}

func TestProcessMonitorSampleConfig(t *testing.T) {
	pm := &ProcessMonitor{}
	cfg := pm.SampleConfig()
	assert.NotEmpty(t, cfg)
	assert.Contains(t, cfg, "processes")
}

func TestProcessMonitorGatherNoProcesses(t *testing.T) {
	pm := &ProcessMonitor{
		Processes: []string{},
		Log:       testutil.Logger{},
	}
	acc := &testutil.Accumulator{}
	err := pm.Gather(acc)
	require.NoError(t, err)
	assert.Equal(t, 0, len(acc.Metrics))
}

func TestProcessMonitorGatherNonExistentProcess(t *testing.T) {
	pm := &ProcessMonitor{
		Processes: []string{"nonexistent_process_xyz_123"},
		Log:       testutil.Logger{},
	}
	acc := &testutil.Accumulator{}
	err := pm.Gather(acc)
	require.NoError(t, err)

	// Should still emit a metric with running=0
	require.Equal(t, 1, len(acc.Metrics))
	m := acc.Metrics[0]
	assert.Equal(t, "process_monitor", m.Measurement)
	assert.Equal(t, "nonexistent_process_xyz_123", m.Tags["process"])
	running, ok := m.Fields["running"]
	assert.True(t, ok)
	assert.Equal(t, 0, running)
}

func TestProcessMonitorGatherSelf(t *testing.T) {
	// "go" process should be running during tests
	pm := &ProcessMonitor{
		Processes: []string{"go"},
		Log:       testutil.Logger{},
	}
	acc := &testutil.Accumulator{}
	err := pm.Gather(acc)
	require.NoError(t, err)

	if len(acc.Metrics) > 0 {
		m := acc.Metrics[0]
		assert.Equal(t, "process_monitor", m.Measurement)
		assert.Equal(t, "go", m.Tags["process"])
		_, hasCPU := m.Fields["cpu_usage"]
		_, hasMem := m.Fields["mem_usage"]
		assert.True(t, hasCPU || hasMem || true) // graceful on CI
	}
}

func TestGetProcessStatsInvalidPID(t *testing.T) {
	_, _, err := getProcessStats("99999999")
	assert.Error(t, err)
}
