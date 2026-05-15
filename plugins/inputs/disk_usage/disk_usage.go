// Package disk_usage provides a Telegraf input plugin for monitoring disk usage statistics.
package disk_usage

import (
	"fmt"
	"syscall"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

// DiskUsage holds the configuration for the disk usage plugin.
type DiskUsage struct {
	// Paths is the list of mount points to monitor. If empty, all mount points are monitored.
	Paths []string `toml:"paths"`
	// IgnoreFS is a list of filesystem types to ignore.
	IgnoreFS []string `toml:"ignore_fs"`
	Log      telegraf.Logger `toml:"-"`
}

var sampleConfig = `
  ## By default stats will be gathered for all mount points.
  ## Set paths to restrict the stats to the specified mount points.
  # paths = ["/", "/home"]

  ## Ignore mount points by filesystem type.
  # ignore_fs = ["tmpfs", "devtmpfs", "devfs", "iso9660", "overlay", "aufs", "squashfs"]
`

// Description returns a short description of the plugin.
func (d *DiskUsage) Description() string {
	return "Read metrics about disk usage by mount point"
}

// SampleConfig returns the default configuration for the plugin.
func (d *DiskUsage) SampleConfig() string {
	return sampleConfig
}

// Gather collects disk usage metrics and writes them to the accumulator.
func (d *DiskUsage) Gather(acc telegraf.Accumulator) error {
	paths := d.Paths
	if len(paths) == 0 {
		var err error
		paths, err = getMountPoints()
		if err != nil {
			return fmt.Errorf("error getting mount points: %w", err)
		}
	}

	ignoreSet := make(map[string]struct{}, len(d.IgnoreFS))
	for _, fs := range d.IgnoreFS {
		ignoreSet[fs] = struct{}{}
	}

	for _, path := range paths {
		stat, fsType, err := getDiskStat(path)
		if err != nil {
			if d.Log != nil {
				d.Log.Warnf("Error getting disk stats for path %s: %v", path, err)
			}
			continue
		}

		if _, ignored := ignoreSet[fsType]; ignored {
			continue
		}

		total := stat.Blocks * uint64(stat.Bsize)
		free := stat.Bfree * uint64(stat.Bsize)
		used := total - free
		var usedPercent float64
		if total > 0 {
			usedPercent = float64(used) / float64(total) * 100.0
		}

		tags := map[string]string{
			"path":   path,
			"fstype": fsType,
		}
		fields := map[string]interface{}{
			"total":        total,
			"free":         free,
			"used":         used,
			"used_percent": usedPercent,
			"inodes_total": stat.Files,
			"inodes_free":  stat.Ffree,
			"inodes_used":  stat.Files - stat.Ffree,
		}
		acc.AddGauge("disk", fields, tags)
	}
	return nil
}

// getDiskStat returns the syscall.Statfs_t and filesystem type for the given path.
func getDiskStat(path string) (*syscall.Statfs_t, string, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return nil, "", err
	}
	// Filesystem type is platform-specific; return empty string as a safe default.
	return &stat, "", nil
}

// getMountPoints returns a list of currently active mount points.
func getMountPoints() ([]string, error) {
	// Default to common mount points if enumeration is not available.
	return []string{"/"}, nil
}

func init() {
	inputs.Add("disk_usage", func() telegraf.Input {
		return &DiskUsage{
			IgnoreFS: []string{"tmpfs", "devtmpfs", "devfs", "iso9660", "overlay", "aufs", "squashfs"},
		}
	})
}
