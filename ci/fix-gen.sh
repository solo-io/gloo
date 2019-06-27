#!/bin/bash

set -e

for file in $(find projects/gloo-ui/src/proto/github.com -type f | grep "_pb.js")
do
    sed -i "s|var gogoproto_gogo_pb = require('../../../../../gogoproto/gogo_pb.js');|var gogoproto_gogo_pb = require('../../../../gogo/protobuf/gogoproto/gogo_pb.js');|g" $file
    printf '%s\n%s\n' "/* eslint-disable */" "$(cat $file)" >$file
done
