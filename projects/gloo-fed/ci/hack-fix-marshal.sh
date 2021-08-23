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
mv projects/gloo-fed/pkg/api/gloo.solo.io/v1/resource_apis.proto projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources.proto > /dev/null
mv projects/gloo-fed/pkg/api/gateway.solo.io/v1/resource_apis.proto projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources.proto > /dev/null
mv projects/gloo-fed/pkg/api/enterprise.gloo.solo.io/v1/resource_apis.proto projects/apiserver/api/rpc.edge.gloo/v1/enterprise_gloo_resources.proto > /dev/null
mv projects/gloo-fed/pkg/api/ratelimit.api.solo.io/v1alpha1/resource_apis.proto projects/apiserver/api/rpc.edge.gloo/v1/ratelimit_resources.proto > /dev/null
mv projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1/resource_apis.proto projects/apiserver/api/fed.rpc/v1/federated_gloo_resources.proto > /dev/null
mv projects/gloo-fed/pkg/api/fed.gateway.solo.io/v1/resource_apis.proto projects/apiserver/api/fed.rpc/v1/federated_gateway_resources.proto > /dev/null
mv projects/gloo-fed/pkg/api/fed.enterprise.gloo.solo.io/v1/resource_apis.proto projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources.proto > /dev/null
mv projects/gloo-fed/pkg/api/fed.ratelimit.solo.io/v1alpha1/resource_apis.proto projects/apiserver/api/fed.rpc/v1/federated_ratelimit_resources.proto > /dev/null

cp projects/apiserver/api/fed.rpc/v1/*resources.proto vendor_any/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/ > /dev/null
cp projects/apiserver/api/rpc.edge.gloo/v1/*resources.proto vendor_any/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/ > /dev/null

# fix cli paths
mkdir -p projects/glooctl-plugins/fed/pkg/api/gloo.solo.io/v1/cli
mv projects/gloo-fed/pkg/api/gloo.solo.io/v1/cli/* projects/glooctl-plugins/fed/pkg/api/gloo.solo.io/v1/cli  > /dev/null
mkdir -p projects/glooctl-plugins/fed/pkg/api/gateway.solo.io/v1/cli
mv projects/gloo-fed/pkg/api/gateway.solo.io/v1/cli/* projects/glooctl-plugins/fed/pkg/api/gateway.solo.io/v1/cli  > /dev/null
mkdir -p projects/glooctl-plugins/fed/pkg/api/enterprise.gloo.solo.io/v1/cli
mv projects/gloo-fed/pkg/api/enterprise.gloo.solo.io/v1/cli/* projects/glooctl-plugins/fed/pkg/api/enterprise.gloo.solo.io/v1/cli  > /dev/null
mkdir -p projects/glooctl-plugins/fed/pkg/api/ratelimit.api.solo.io/v1alpha1/cli
mv projects/gloo-fed/pkg/api/ratelimit.api.solo.io/v1alpha1/cli/* projects/glooctl-plugins/fed/pkg/api/ratelimit.api.solo.io/v1alpha1/cli  > /dev/null


# All of the generated files get created under projects/gloo-fed by default. Move the apiserver-related files to the appropriate apiserver folders:

# Create folders for the base resource handlers
mkdir -p projects/apiserver/pkg/api/gloo.solo.io/v1/handler
mkdir -p projects/apiserver/pkg/api/gateway.solo.io/v1/handler
mkdir -p projects/apiserver/pkg/api/enterprise.gloo.solo.io/v1/handler
mkdir -p projects/apiserver/pkg/api/ratelimit.api.solo.io/v1alpha1/handler

# Move Gloo Fed base resource handlers into their pkg dirs
mv projects/gloo-fed/pkg/api/gloo.solo.io/v1/handler/fed_handler.go projects/apiserver/pkg/api/gloo.solo.io/v1/handler  > /dev/null
mv projects/gloo-fed/pkg/api/gateway.solo.io/v1/handler/fed_handler.go projects/apiserver/pkg/api/gateway.solo.io/v1/handler  > /dev/null
mv projects/gloo-fed/pkg/api/enterprise.gloo.solo.io/v1/handler/fed_handler.go projects/apiserver/pkg/api/enterprise.gloo.solo.io/v1/handler  > /dev/null
mv projects/gloo-fed/pkg/api/ratelimit.api.solo.io/v1alpha1/handler/fed_handler.go projects/apiserver/pkg/api/ratelimit.api.solo.io/v1alpha1/handler  > /dev/null

# Move the single-cluster base resource checkers (i.e. the files that return resource summaries for GlooInstances) into their pkg dirs
mv projects/gloo-fed/pkg/api/gloo.solo.io/v1/check/single_cluster_check.go projects/apiserver/pkg/api/gloo.solo.io/v1/handler  > /dev/null
mv projects/gloo-fed/pkg/api/gateway.solo.io/v1/check/single_cluster_check.go projects/apiserver/pkg/api/gateway.solo.io/v1/handler  > /dev/null
mv projects/gloo-fed/pkg/api/enterprise.gloo.solo.io/v1/check/single_cluster_check.go projects/apiserver/pkg/api/enterprise.gloo.solo.io/v1/handler  > /dev/null
mv projects/gloo-fed/pkg/api/ratelimit.api.solo.io/v1alpha1/check/single_cluster_check.go projects/apiserver/pkg/api/ratelimit.api.solo.io/v1alpha1/handler  > /dev/null

# The single-cluster base resource handlers depend on the glooinstance handler, which in turn depends on the pkg/api handlers above, so to avoid circular dependencies,
# we put the base resource handlers in a different folder instead of the pkg/api folders above
mv projects/gloo-fed/pkg/api/gloo.solo.io/v1/handler/single_cluster_handler.go projects/apiserver/server/services/single_cluster_resource_handler/single_cluster_gloo_handler.go  > /dev/null
mv projects/gloo-fed/pkg/api/gateway.solo.io/v1/handler/single_cluster_handler.go projects/apiserver/server/services/single_cluster_resource_handler/single_cluster_gateway_handler.go  > /dev/null
mv projects/gloo-fed/pkg/api/enterprise.gloo.solo.io/v1/handler/single_cluster_handler.go projects/apiserver/server/services/single_cluster_resource_handler/single_cluster_enterprise_gloo_handler.go  > /dev/null
mv projects/gloo-fed/pkg/api/ratelimit.api.solo.io/v1alpha1/handler/single_cluster_handler.go projects/apiserver/server/services/single_cluster_resource_handler/single_cluster_ratelimit_handler.go  > /dev/null

# Create folders for the federated resource handlers and move the files over.
mkdir -p projects/apiserver/pkg/api/fed.gloo.solo.io/v1/handler
mv projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1/handler/* projects/apiserver/pkg/api/fed.gloo.solo.io/v1/handler  > /dev/null
mkdir -p projects/apiserver/pkg/api/fed.gateway.solo.io/v1/handler
mv projects/gloo-fed/pkg/api/fed.gateway.solo.io/v1/handler/* projects/apiserver/pkg/api/fed.gateway.solo.io/v1/handler  > /dev/null
mkdir -p projects/apiserver/pkg/api/fed.enterprise.gloo.solo.io/v1/handler
mv projects/gloo-fed/pkg/api/fed.enterprise.gloo.solo.io/v1/handler/* projects/apiserver/pkg/api/fed.enterprise.gloo.solo.io/v1/handler  > /dev/null
mkdir -p projects/apiserver/pkg/api/fed.ratelimit.solo.io/v1alpha1/handler
mv projects/gloo-fed/pkg/api/fed.ratelimit.solo.io/v1alpha1/handler/* projects/apiserver/pkg/api/fed.ratelimit.solo.io/v1alpha1/handler  > /dev/null
