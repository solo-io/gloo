#!/bin/bash

set -e

# protoc generation for javascript is broken
# when protos import other protos, the import path relative directories are often broken
# this script exists to fix such compilation errors

# to test if ui compiles:
# yarn tsc --noEmit (quicker, doesn't catch all errors)
# yarn build        (slower, is run in CI will catch everything)

for file in $(find projects/gloo-ui/src/proto -type f | grep "_pb.js")
do
  sed "s|google/api/annotations_pb.js|solo-kit/api/external/google/api/annotations_pb.js|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|google/api/http_pb.js|solo-kit/api/external/google/api/http_pb.js|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|envoy/api/v2/core/http_uri_pb.js|solo-kit/api/external/envoy/api/v2/core/http_uri_pb.js|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|google/rpc/status_pb.js|solo-kit/api/external/google/rpc/status_pb.js|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|envoy/api/v2/discovery_pb.js|solo-kit/api/external/envoy/api/v2/discovery_pb.js|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|envoy/api/v2/core/base_pb.js|solo-kit/api/external/envoy/api/v2/core/base_pb.js|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|envoy/type/range_pb.js|gloo/projects/gloo/api/external/envoy/type/range_pb.js|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|envoy/type/percent_pb.js|solo-kit/api/external/envoy/type/percent_pb.js|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|envoy/api/v2/route/route_pb.js|gloo/projects/gloo/api/external/envoy/api/v2/route/route_pb.js|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|../../../../../extproto/ext_pb.js|../../../../../protoc-gen-ext/extproto/ext_pb.js|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|../../../extproto/ext_pb.js|../../../protoc-gen-ext/extproto/ext_pb.js|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  printf '%s\n%s\n' "/* eslint-disable */" "$(cat "$file")" > "$file"
done
