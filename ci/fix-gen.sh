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
  sed "s|google/api/annotations_pb.js|github.com/solo-io/solo-kit/api/external/google/api/annotations_pb.js|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|google/api/http_pb.js|github.com/solo-io/solo-kit/api/external/google/api/http_pb.js|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|google/rpc/status_pb.js|github.com/solo-io/solo-kit/api/external/google/rpc/status_pb.js|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|envoy/annotations/|github.com/solo-io/gloo/projects/gloo/api/external/envoy/annotations/|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|udpa/annotations/|github.com/solo-io/gloo/projects/gloo/api/external/udpa/annotations/|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|../../../../../protoc-gen-ext/extproto/ext_pb.js|../../../../../extproto/ext_pb.js |g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  printf '%s\n%s\n' "/* eslint-disable */" "$(cat "$file")" > "$file"
done

for file in $(find projects/gloo-ui/src/proto -type f | grep "_pb.d.ts")
do
  sed "s|../../../../extproto/ext_pb|../../../../protoc-gen-ext/extproto/ext_pb|g" "$file" > "$file".tmp && mv "$file".tmp "$file"

  printf '%s\n%s\n' "/* eslint-disable */" "$(cat "$file")" > "$file"
done
