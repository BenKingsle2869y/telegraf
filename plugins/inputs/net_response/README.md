# Net Response Input Plugin

The `net_response` plugin checks network connectivity to a given address and port.
It measures the response time and reports whether the connection succeeded.

## Configuration

```toml
[[inputs.net_response]]
  ## Protocol to check (tcp or udp)
  protocol = "tcp"

  ## Address and port to check
  address = "localhost:80"

  ## Timeout for the connection
  ## Increase this if you're monitoring hosts across a slow or high-latency network
  timeout = "2s"
```

## Metrics

- `net_response`
  - tags:
    - `protocol` — the protocol used (`tcp` or `udp`)
    - `address` — the target address and port
  - fields:
    - `response_time` (float) — time in seconds to establish the connection
    - `result_code` (int) — `0` for success, `1` for failure
    - `result` (string) — `"success"` or an error message

## Example Output

```
net_response,address=localhost:80,protocol=tcp response_time=0.001234,result="success",result_code=0i
net_response,address=localhost:9999,protocol=tcp response_time=0.500123,result="error: connection refused",result_code=1i
```
