# Process Monitor Input Plugin

The `process_monitor` plugin monitors CPU and memory usage of specific named
processes on the host system. It uses `pgrep` and `ps` to collect per-process
metrics.

## Requirements

- Linux / macOS (requires `pgrep` and `ps` commands)
- Telegraf running with sufficient permissions to inspect target processes

## Configuration

```toml
[[inputs.process_monitor]]
  ## List of process names to monitor
  processes = ["telegraf", "nginx", "postgres"]
```

## Metrics

Measurement: `process_monitor`

### Tags

| Tag       | Description                        |
|-----------|------------------------------------|
| `process` | Name of the monitored process      |

### Fields

| Field       | Type  | Description                                      |
|-------------|-------|--------------------------------------------------|
| `running`   | int   | 1 if process is running, 0 otherwise             |
| `pid_count` | int   | Number of PIDs found for the process             |
| `cpu_usage` | float | Total CPU usage (%) across all matching PIDs     |
| `mem_usage` | float | Total memory usage (%) across all matching PIDs  |

## Example Output

```
process_monitor,process=nginx running=1i,pid_count=4i,cpu_usage=0.4,mem_usage=1.2 1700000000000000000
process_monitor,process=postgres running=1i,pid_count=2i,cpu_usage=0.1,mem_usage=2.5 1700000000000000000
process_monitor,process=unknown running=0i 1700000000000000000
```
