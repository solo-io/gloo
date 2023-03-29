#!/usr/bin/env python3
"""
This is a quick script to take in in a .csv file and output a .csv file with
the same data, but with the following format:
    - 1st column is unchanged
    - 2nd column is time in seconds since the first data point
"""

from datetime import datetime
import sys

infile=sys.argv[1]
outfile=sys.argv[2]

out = open(outfile, "a")
out.write("number_of_upstreams,timestamp\n")

original = -1
for line in open(infile, 'r').readlines():
    line = line.strip()
    num_upstreams, t = line.split(",")

    if "number_of_upstreams" in line:
        continue
    if original == -1:
        original = datetime.strptime(t, "%Y-%m-%d %H:%M:%S")
    current = datetime.strptime(t, "%Y-%m-%d %H:%M:%S")

    t = int((current - original).total_seconds())

    out.write("%s,%s\n" % (num_upstreams, t))

out.close()