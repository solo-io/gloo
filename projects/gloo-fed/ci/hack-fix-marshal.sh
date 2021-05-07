#!/bin/bash

set -e

for file in $(find projects/gloo-fed/pkg/api/fed.solo.io/v1/input -type f | grep ".go")
do
  sed "s|\"github.com/solo-io/solo-apis/projects/gloo-fed/pkg/api/gateway.solo.io/v1|\"github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|\"github.com/solo-io/solo-apis/projects/gloo-fed/pkg/api/gloo.solo.io/v1|\"github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|\"github.com/solo-io/solo-apis/projects/gloo-fed/pkg/api/enterprise.gloo.solo.io/v1|\"github.com/solo-io/solo-apis/pkg/api/enterprise.gloo.solo.io/v1|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|\"github.com/solo-io/solo-apis/projects/gloo-fed/pkg/api/ratelimit.api.solo.io/v1alpha1|\"github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|\"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/gateway.solo.io/v1|\"github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|\"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/gloo.solo.io/v1|\"github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|\"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/enterprise.gloo.solo.io/v1|\"github.com/solo-io/solo-apis/pkg/api/enterprise.gloo.solo.io/v1|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
  sed "s|\"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/ratelimit.api.solo.io/v1alpha1|\"github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
done

# fix the fact that gloo doesn't have an upstream_group.proto file
for file in projects/gloo-fed/pkg/api/gloo.solo.io/v1/resource_apis.proto
do
  sed "s|import \"github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_group.proto\";||g" "$file" > "$file".tmp && mv "$file".tmp "$file"
done

# fix the ratelimit_resources.proto file
for file in projects/gloo-fed/pkg/api/ratelimit.api.solo.io/v1alpha1/resource_apis.proto projects/gloo-fed/pkg/api/fed.ratelimit.solo.io/v1alpha1/resource_apis.proto
do
  sed "s|github.com/solo-io/solo-apis/api/gloo/ratelimit/v1alpha1/rate_limit_config.proto|github.com/solo-io/solo-apis/api/rate-limiter/v1alpha1/ratelimit.proto|g" "$file" > "$file".tmp && mv "$file".tmp "$file"
done

# fix generated proto paths
mv projects/gloo-fed/pkg/api/gloo.solo.io/v1/resource_apis.proto projects/apiserver/api/fed.rpc/v1/gloo_resources.proto > /dev/null
mv projects/gloo-fed/pkg/api/gateway.solo.io/v1/resource_apis.proto projects/apiserver/api/fed.rpc/v1/gateway_resources.proto > /dev/null
mv projects/gloo-fed/pkg/api/enterprise.gloo.solo.io/v1/resource_apis.proto projects/apiserver/api/fed.rpc/v1/enterprise_gloo_resources.proto > /dev/null
mv projects/gloo-fed/pkg/api/ratelimit.api.solo.io/v1alpha1/resource_apis.proto projects/apiserver/api/fed.rpc/v1/ratelimit_resources.proto > /dev/null
mv projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1/resource_apis.proto projects/apiserver/api/fed.rpc/v1/federated_gloo_resources.proto > /dev/null
mv projects/gloo-fed/pkg/api/fed.gateway.solo.io/v1/resource_apis.proto projects/apiserver/api/fed.rpc/v1/federated_gateway_resources.proto > /dev/null
mv projects/gloo-fed/pkg/api/fed.enterprise.gloo.solo.io/v1/resource_apis.proto projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources.proto > /dev/null
mv projects/gloo-fed/pkg/api/fed.ratelimit.solo.io/v1alpha1/resource_apis.proto projects/apiserver/api/fed.rpc/v1/federated_ratelimit_resources.proto > /dev/null
cp projects/apiserver/api/fed.rpc/v1/*resources.proto vendor_any/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/ > /dev/null

# fix cli paths
mkdir -p projects/glooctl-plugins/fed/pkg/api/gloo.solo.io/v1/cli
mv projects/gloo-fed/pkg/api/gloo.solo.io/v1/cli/* projects/glooctl-plugins/fed/pkg/api/gloo.solo.io/v1/cli  > /dev/null
mkdir -p projects/glooctl-plugins/fed/pkg/api/gateway.solo.io/v1/cli
mv projects/gloo-fed/pkg/api/gateway.solo.io/v1/cli/* projects/glooctl-plugins/fed/pkg/api/gateway.solo.io/v1/cli  > /dev/null
mkdir -p projects/glooctl-plugins/fed/pkg/api/enterprise.gloo.solo.io/v1/cli
mv projects/gloo-fed/pkg/api/enterprise.gloo.solo.io/v1/cli/* projects/glooctl-plugins/fed/pkg/api/enterprise.gloo.solo.io/v1/cli  > /dev/null
mkdir -p projects/glooctl-plugins/fed/pkg/api/ratelimit.api.solo.io/v1alpha1/cli
mv projects/gloo-fed/pkg/api/ratelimit.api.solo.io/v1alpha1/cli/* projects/glooctl-plugins/fed/pkg/api/ratelimit.api.solo.io/v1alpha1/cli  > /dev/null

# fix apiserver handler paths
mkdir -p projects/apiserver/pkg/api/gloo.solo.io/v1/handler
mv projects/gloo-fed/pkg/api/gloo.solo.io/v1/handler/* projects/apiserver/pkg/api/gloo.solo.io/v1/handler  > /dev/null
mkdir -p projects/apiserver/pkg/api/gateway.solo.io/v1/handler
mv projects/gloo-fed/pkg/api/gateway.solo.io/v1/handler/* projects/apiserver/pkg/api/gateway.solo.io/v1/handler  > /dev/null
mkdir -p projects/apiserver/pkg/api/enterprise.gloo.solo.io/v1/handler
mv projects/gloo-fed/pkg/api/enterprise.gloo.solo.io/v1/handler/* projects/apiserver/pkg/api/enterprise.gloo.solo.io/v1/handler  > /dev/null
mkdir -p projects/apiserver/pkg/api/ratelimit.api.solo.io/v1alpha1/handler
mv projects/gloo-fed/pkg/api/ratelimit.api.solo.io/v1alpha1/handler/* projects/apiserver/pkg/api/ratelimit.api.solo.io/v1alpha1/handler  > /dev/null

mkdir -p projects/apiserver/pkg/api/fed.gloo.solo.io/v1/handler
mv projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1/handler/* projects/apiserver/pkg/api/fed.gloo.solo.io/v1/handler  > /dev/null
mkdir -p projects/apiserver/pkg/api/fed.gateway.solo.io/v1/handler
mv projects/gloo-fed/pkg/api/fed.gateway.solo.io/v1/handler/* projects/apiserver/pkg/api/fed.gateway.solo.io/v1/handler  > /dev/null
mkdir -p projects/apiserver/pkg/api/fed.enterprise.gloo.solo.io/v1/handler
mv projects/gloo-fed/pkg/api/fed.enterprise.gloo.solo.io/v1/handler/* projects/apiserver/pkg/api/fed.enterprise.gloo.solo.io/v1/handler  > /dev/null
mkdir -p projects/apiserver/pkg/api/fed.ratelimit.solo.io/v1alpha1/handler
mv projects/gloo-fed/pkg/api/fed.ratelimit.solo.io/v1alpha1/handler/* projects/apiserver/pkg/api/fed.ratelimit.solo.io/v1alpha1/handler  > /dev/null
