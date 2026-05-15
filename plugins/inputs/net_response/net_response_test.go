package net_response

import (
	"net"
	"testing"
	"time"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNetResponseDescription(t *testing.T) {
	n := &NetResponse{}
	assert.Equal(t, description, n.Description())
}

func TestNetResponseSampleConfig(t *testing.T) {
	n := &NetResponse{}
	assert.Equal(t, sampleConfig, n.SampleConfig())
}

func TestNetResponseDefaults(t *testing.T) {
	n := &NetResponse{}
	acc := &testutil.Accumulator{}
	// Should not panic with empty config; defaults applied in Gather
	_ = n.Gather(acc)
	assert.Equal(t, "tcp", n.Protocol)
	assert.Equal(t, time.Second, n.Timeout)
}

func TestNetResponseTCPSuccess(t *testing.T) {
	// Start a local TCP listener
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()

	n := &NetResponse{
		Protocol: "tcp",
		Address:  ln.Addr().String(),
		Timeout:  time.Second,
	}

	acc := &testutil.Accumulator{}
	err = n.Gather(acc)
	require.NoError(t, err)

	acc.AssertContainsTaggedFields(t, "net_response",
		map[string]interface{}{
			"result_code": 0,
			"result":      "success",
		},
		map[string]string{
			"protocol": "tcp",
			"address":  ln.Addr().String(),
		},
	)
}

func TestNetResponseTCPFailure(t *testing.T) {
	n := &NetResponse{
		Protocol: "tcp",
		Address:  "127.0.0.1:19999",
		Timeout:  500 * time.Millisecond,
	}

	acc := &testutil.Accumulator{}
	err := n.Gather(acc)
	require.NoError(t, err)

	fields, ok := acc.Get("net_response")
	require.True(t, ok)
	assert.Equal(t, 1, fields.Fields["result_code"])
}
