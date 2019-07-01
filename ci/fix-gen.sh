#!/bin/bash

set -e

for file in $(find projects/gloo-ui/src/proto/github.com -type f | grep "_pb.js")
do
	sed -i "s|validate/validate_pb.js|github.com/solo-io/solo-kit/api/external/validate/validate_pb.js|g" $file
	sed -i "s|google/api/annotations_pb.js|github.com/solo-io/solo-kit/api/external/google/api/annotations_pb.js|g" $file
	sed -i "s|google/api/http_pb.js|github.com/solo-io/solo-kit/api/external/google/api/http_pb.js|g" $file
	sed -i "s|google/rpc/status_pb.js|github.com/solo-io/solo-kit/api/external/google/rpc/status_pb.js|g" $file
	sed -i "s|envoy/api/v2/discovery_pb.js|github.com/solo-io/solo-kit/api/external/envoy/api/v2/discovery_pb.js|g" $file
	sed -i "s|envoy/api/v2/core/base_pb.js|github.com/solo-io/solo-kit/api/external/envoy/api/v2/core/base_pb.js|g" $file
	sed -i "s|envoy/type/percent_pb.js|github.com/solo-io/solo-kit/api/external/envoy/type/percent_pb.js|g" $file
    sed -i "s|../../../../../gogoproto/gogo_pb.js|../../../../gogo/protobuf/gogoproto/gogo_pb.js|g" $file
    printf '%s\n%s\n' "/* eslint-disable */" "$(cat $file)" >$file
done
