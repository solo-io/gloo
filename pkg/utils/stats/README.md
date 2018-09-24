
# Logs

enable debug logs with 
```shell
curl -XPUT -d '{"level":"debug"}' -H"content-type: application/json" http://localhost:9091/logging
```

# zPages

see them here:

localhost:9091/zpages/tracez

# Local Prometheus

Start like so:
```shell
cat <<EOF > config.yaml
global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.

# A scrape configuration containing exactly one endpoint to scrape:
scrape_configs:
  # The job name is added as a label 'job=<job_name>' to any timeseries scraped from this config.
  - job_name: 'solo-kit'

    # Override the global default and scrape targets from this job every 1 seconds.
    scrape_interval: 1s

    static_configs:
      - targets: ['localhost:9091']
EOF
prometheus --config.file=./config.yaml
```

Prometheus will be available in `localhost:9090`.

To see the rate of incoming snapshots, try this query:
```
rate(api_snap_emitter_snap_in[5m])
```

# Profiling
## CPU tracing
```
$ curl http://localhost:9091/debug/pprof/trace?seconds=5 -o trace.out
$ go tool trace trace.out
```

## profiling
```
$ curl http://localhost:9091/debug/pprof/profile?seconds=5 -o pprof.out
$ go tool pprof pprof.out
(pprof) top
(pprof) web
```

## Stack traces
See go routines:
```
$ curl http://localhost:9091/debug/pprof/goroutine?debug=2
```