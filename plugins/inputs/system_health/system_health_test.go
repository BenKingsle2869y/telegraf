package system_health

import (
	"testing"
	"time"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemHealthDefaults(t *testing.T) {
	s := &SystemHealth{}
	assert.False(t, s.CollectCPU)
	assert.False(t, s.CollectMemory)
	assert.Equal(t, 0.0, s.CPUWarning)
	assert.Equal(t, 0.0, s.MemWarning)
}

func TestSystemHealthDescription(t *testing.T) {
	s := &SystemHealth{}
	assert.NotEmpty(t, s.Description())
}

func TestSystemHealthSampleConfig(t *testing.T) {
	s := &SystemHealth{}
	assert.NotEmpty(t, s.SampleConfig())
}

func TestGatherCPUOnly(t *testing.T) {
	s := &SystemHealth{
		CollectCPU:    true,
		CollectMemory: false,
		CPUWarning:    90.0,
		MemWarning:    90.0,
	}

	acc := &testutil.Accumulator{}
	err := s.Gather(acc)
	require.NoError(t, err)

	acc.Wait(1)
	assert.True(t, acc.HasMeasurement("system_health_cpu"))
	assert.False(t, acc.HasMeasurement("system_health_memory"))

	metric, ok := acc.Get("system_health_cpu")
	require.True(t, ok)
	assert.Contains(t, metric.Fields, "usage_percent")
	assert.Contains(t, metric.Fields, "warning")
	assert.WithinDuration(t, time.Now(), metric.Time, 5*time.Second)
}

func TestGatherMemoryOnly(t *testing.T) {
	s := &SystemHealth{
		CollectCPU:    false,
		CollectMemory: true,
		CPUWarning:    90.0,
		MemWarning:    90.0,
	}

	acc := &testutil.Accumulator{}
	err := s.Gather(acc)
	require.NoError(t, err)

	acc.Wait(1)
	assert.False(t, acc.HasMeasurement("system_health_cpu"))
	assert.True(t, acc.HasMeasurement("system_health_memory"))

	metric, ok := acc.Get("system_health_memory")
	require.True(t, ok)
	assert.Contains(t, metric.Fields, "used_percent")
	assert.Contains(t, metric.Fields, "total_bytes")
	assert.Contains(t, metric.Fields, "used_bytes")
	assert.Contains(t, metric.Fields, "free_bytes")
	assert.Contains(t, metric.Fields, "warning")
}

func TestGatherBothDisabled(t *testing.T) {
	s := &SystemHealth{
		CollectCPU:    false,
		CollectMemory: false,
	}

	acc := &testutil.Accumulator{}
	err := s.Gather(acc)
	require.NoError(t, err)
	assert.Equal(t, 0, len(acc.Metrics))
}
