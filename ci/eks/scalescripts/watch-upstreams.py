#!/usr/bin/env python3
"""
This is a quick script to watch the progress of a given gloo federation operation
"""

import subprocess
import time
import sys

kube_context = sys.argv[1]
outfile = sys.argv[2]
iterations = int(sys.argv[3])

# lazy access to bash shell
def shell(cmd):
    process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
    output, error = process.communicate()
    return output

# quick kube command to query all upstreams from kube_context
def count_upstreams():
    cmd = "kubectl --context=" + kube_context +  ' --namespace=gloo-system get upstream --no-headers --selector=fed.solo.io/owner'
    out = shell(cmd)
    return len(out.decode("utf-8").split('\n'))

# quickly create a simple csv

f = open(outfile, "a")
f.write("number_of_upstreams,timestamp\n")
f.close()

upstream_count = 0

while upstream_count < iterations * 3:
    # fetch data every 2 seconds (plus the delay of the actual command)
    time.sleep(2)
    upstream_count = count_upstreams()
    print("Total upstreams: " + str(upstream_count))
    to_write = str(upstream_count) + "," + time.strftime("%Y-%m-%d %H:%M:%S", time.gmtime())

    # write data to file (open and close each time, so that writes are visible)
    f = open(outfile, "a")
    f.write(to_write + "\n")
    f.close()

    # allow terminal as another source of truth/progression
    print(to_write)
