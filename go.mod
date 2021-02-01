module github.com/solo-io/solo-projects

go 1.14

require (
	cloud.google.com/go/datastore v1.3.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/avast/retry-go v2.4.3+incompatible
	github.com/aws/aws-sdk-go v1.34.9
	github.com/deislabs/oras v0.8.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/elazarl/goproxy/ext v0.0.0-20191011121108-aa519ddbe484 // indirect
	github.com/envoyproxy/go-control-plane v0.9.8
	github.com/envoyproxy/protoc-gen-validate v0.4.0
	github.com/fgrosse/zaptest v1.1.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-logr/zapr v0.2.0
	github.com/go-redis/redis/v8 v8.2.3
	github.com/go-test/deep v1.0.4
	github.com/gobuffalo/logger v1.0.3 // indirect
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.3
	github.com/google/wire v0.4.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/hashicorp/consul/api v1.3.0
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/vault/api v1.0.5-0.20191108163347-bdd38fca2cff
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keybase/go-ps v0.0.0-20190827175125-91aafc93ba19
	github.com/mitchellh/hashstructure v1.0.0
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/opencontainers/go-digest v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.8.0
	github.com/prometheus/prometheus v2.5.0+incompatible
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/radovskyb/watcher v1.0.7 // indirect
	github.com/rotisserie/eris v0.4.0
	github.com/sirupsen/logrus v1.7.0
	github.com/solo-io/anyvendor v0.0.1
	github.com/solo-io/envoy-operator v0.1.4
	github.com/solo-io/ext-auth-plugins v0.2.1
	github.com/solo-io/ext-auth-service v0.7.10
	github.com/solo-io/external-apis v0.1.4
	github.com/solo-io/gloo v1.7.0-beta12
	github.com/solo-io/go-utils v0.20.2
	github.com/solo-io/k8s-utils v0.0.5
	github.com/solo-io/licensing v0.1.20
	github.com/solo-io/protoc-gen-ext v0.0.15
	github.com/solo-io/rate-limiter v0.1.9
	github.com/solo-io/reporting-client v0.2.0
	github.com/solo-io/skv2 v0.17.0
	github.com/solo-io/skv2-enterprise v0.0.8
	github.com/solo-io/solo-apis v0.0.0-20210126142406-19188dc1c355
	github.com/solo-io/solo-kit v0.17.1
	github.com/solo-io/wasm-kit v0.1.3
	github.com/solo-io/wasm/tools/wasme/pkg v0.0.0-20201021213306-77f82bdc3cc3
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/tredoe/osutil v0.0.0-20191018075336-e272fdda81c8 // indirect
	go.opencensus.io v0.22.5
	go.uber.org/zap v1.16.0
	golang.org/x/mod v0.3.0
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b
	golang.org/x/tools v0.0.0-20210106214847-113979e3529a
	google.golang.org/genproto v0.0.0-20201019141844-1ed22bb0c154
	google.golang.org/grpc v1.33.1
	google.golang.org/protobuf v1.25.0
	gopkg.in/square/go-jose.v2 v2.5.1
	helm.sh/helm/v3 v3.4.2
	k8s.io/api v0.19.7
	k8s.io/apiextensions-apiserver v0.19.7
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/code-generator v0.19.7
	k8s.io/kubernetes v1.19.6
	k8s.io/utils v0.19.6 // indirect
	sigs.k8s.io/controller-runtime v0.7.0
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.4.3
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

	github.com/golang/mock v1.4.4-0.20200406172829-6d816de489c1 => github.com/golang/mock v1.4.3
	// needed by gloo
	github.com/google/go-github/v32 => github.com/google/go-github/v32 v32.0.0
	github.com/opencontainers/go-digest => github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/sclevine/agouti => github.com/yuval-k/agouti v0.0.0-20190109124522-0e71d6bad483
	// Lock sys package to fix darwin upgrade issue
	// https://github.com/helm/chart-releaser/pull/82/files#diff-33ef32bf6c23acb95f5902d7097b7a1d5128ca061167ec0716715b0b9eeaa5f6R41
	golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
	k8s.io/api => k8s.io/api v0.19.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.6
	k8s.io/apiserver => k8s.io/apiserver v0.19.6
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.6
	k8s.io/client-go => k8s.io/client-go v0.19.6
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.19.6
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.19.6
	k8s.io/code-generator => k8s.io/code-generator v0.19.6
	k8s.io/component-base => k8s.io/component-base v0.19.6
	k8s.io/cri-api => k8s.io/cri-api v0.19.6
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.19.6
	k8s.io/gengo => k8s.io/gengo v0.0.0-20190822140433-26a664648505
	k8s.io/heapster => k8s.io/heapster v1.2.0-beta.1
	k8s.io/klog => github.com/stefanprodan/klog v0.0.0-20190418165334-9cbb78b20423
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.19.6
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.19.6
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.19.6
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.19.6
	k8s.io/kubectl => k8s.io/kubectl v0.19.6
	k8s.io/kubelet => k8s.io/kubelet v0.19.6
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.19.6
	k8s.io/metrics => k8s.io/metrics v0.19.6
	k8s.io/repo-infra => k8s.io/repo-infra v0.0.0-20181204233714-00fe14e3d1a3
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.19.6
	k8s.io/sample-controller => k8s.io/sample-controller v0.19.6
	k8s.io/utils => k8s.io/utils v0.0.0-20200729134348-d5654de09c73
)
