package net_response

import (
	"fmt"
	"net"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

const (
	sampleConfig = `
  ## Protocol to check (tcp or udp)
  protocol = "tcp"

  ## Address and port to check
  address = "localhost:80"

  ## Timeout for the connection
  ## Increase this if you're monitoring slow or remote endpoints
  timeout = "3s"
`
	description = "Checks network connectivity to a given address and port"
)

// NetResponse defines the plugin structure
type NetResponse struct {
	Protocol string        `toml:"protocol"`
	Address  string        `toml:"address"`
	Timeout  time.Duration `toml:"timeout"`
}

func (n *NetResponse) Description() string {
	return description
}

func (n *NetResponse) SampleConfig() string {
	return sampleConfig
}

func (n *NetResponse) Gather(acc telegraf.Accumulator) error {
	if n.Protocol == "" {
		n.Protocol = "tcp"
	}
	// Default timeout bumped to 3s to reduce false positives on slower networks
	if n.Timeout == 0 {
		n.Timeout = 3 * time.Second
	}

	start := time.Now()
	conn, err := net.DialTimeout(n.Protocol, n.Address, n.Timeout)
	elapsed := time.Since(start).Seconds()

	tags := map[string]string{
		"protocol": n.Protocol,
		"address":  n.Address,
	}

	fields := map[string]interface{}{
		"response_time": elapsed,
	}

	if err != nil {
		fields["result_code"] = 1
		fields["result"] = fmt.Sprintf("error: %s", err.Error())
	} else {
		conn.Close()
		fields["result_code"] = 0
		fields["result"] = "success"
	}

	acc.AddFields("net_response", fields, tags)
	return nil
}

func init() {
	inputs.Add("net_response", func() telegraf.Input {
		return &NetResponse{
			Protocol: "tcp",
			Timeout:  3 * time.Second,
		}
	})
}
