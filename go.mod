module github.com/solo-io/solo-projects

go 1.14

require (
	github.com/avast/retry-go v2.4.3+incompatible
	github.com/aws/aws-sdk-go v1.30.15
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/elazarl/goproxy/ext v0.0.0-20191011121108-aa519ddbe484 // indirect
	github.com/envoyproxy/go-control-plane v0.9.6-0.20200529035633-fc42e08917e9
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/fgrosse/zaptest v1.1.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/gogo/googleapis v1.3.1
	github.com/gogo/protobuf v1.3.1
	github.com/golang/mock v1.4.4-0.20200406172829-6d816de489c1
	github.com/golang/protobuf v1.3.5
	github.com/google/wire v0.4.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/hashicorp/consul/api v1.3.0
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/vault/api v1.0.5-0.20191108163347-bdd38fca2cff
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keybase/go-ps v0.0.0-20190827175125-91aafc93ba19
	github.com/mitchellh/hashstructure v1.0.0
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.2.1
	github.com/prometheus/prometheus v2.5.0+incompatible
	github.com/rotisserie/eris v0.1.1
	github.com/solo-io/anyvendor v0.0.1
	github.com/solo-io/dev-portal v0.1.12
	github.com/solo-io/envoy-operator v0.1.1
	github.com/solo-io/ext-auth-plugin-examples v0.1.2-0.20191219182743-362e50211779
	github.com/solo-io/ext-auth-plugins v0.1.2
	github.com/solo-io/ext-auth-service v0.6.13
	github.com/solo-io/gloo v1.5.0-beta8
	github.com/solo-io/go-utils v0.16.4
	github.com/solo-io/licensing v0.1.17
	github.com/solo-io/protoc-gen-ext v0.0.9
	github.com/solo-io/rate-limiter v0.0.3
	github.com/solo-io/reporting-client v0.1.2
	github.com/solo-io/solo-apis v0.0.0-20200709192249-4f741634437b
	github.com/solo-io/solo-kit v0.13.9
	github.com/tredoe/osutil v0.0.0-20191018075336-e272fdda81c8 // indirect
	go.opencensus.io v0.22.2
	go.uber.org/atomic v1.5.0
	go.uber.org/zap v1.13.0
	golang.org/x/mod v0.3.0
	golang.org/x/net v0.0.0-20200625001655-4c5254603344
	golang.org/x/tools v0.0.0-20200717024301-6ddee64345a6
	google.golang.org/genproto v0.0.0-20200309141739-5b75447e413d
	google.golang.org/grpc v1.28.0-pre.0.20200226185027-6cd03861bfd2
	gopkg.in/square/go-jose.v2 v2.3.1
	helm.sh/helm/v3 v3.1.2
	k8s.io/api v0.17.2
	k8s.io/apiextensions-apiserver v0.18.1
	k8s.io/apimachinery v0.18.1
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/code-generator v0.17.2
	k8s.io/kubernetes v1.17.1
	sigs.k8s.io/controller-runtime v0.5.6
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.0.0+incompatible
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.0.5
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

	github.com/golang/mock v1.4.4-0.20200406172829-6d816de489c1 => github.com/golang/mock v1.4.3
	github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1
	github.com/hashicorp/consul/api => github.com/hashicorp/consul/api v1.1.0
	github.com/hashicorp/vault => github.com/hashicorp/vault v1.3.0
	github.com/hashicorp/vault/api => github.com/hashicorp/vault/api v1.0.5-0.20191108163347-bdd38fca2cff
	// We need to pin this to the version used by solo-kit, because dev-portal currently uses v1.3.1 (incompatible)
	github.com/pseudomuto/protoc-gen-doc => github.com/pseudomuto/protoc-gen-doc v1.0.0
	github.com/sclevine/agouti => github.com/yuval-k/agouti v0.0.0-20190109124522-0e71d6bad483

	// kube 1.17
	k8s.io/api => k8s.io/api v0.17.1
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.1
	k8s.io/apiserver => k8s.io/apiserver v0.17.1
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.17.1
	k8s.io/client-go => k8s.io/client-go v0.17.1
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.17.1
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.1
	k8s.io/code-generator => k8s.io/code-generator v0.17.1
	k8s.io/component-base => k8s.io/component-base v0.17.1
	k8s.io/cri-api => k8s.io/cri-api v0.17.1
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.17.1
	k8s.io/gengo => k8s.io/gengo v0.0.0-20190822140433-26a664648505
	k8s.io/heapster => k8s.io/heapster v1.2.0-beta.1
	k8s.io/klog => github.com/stefanprodan/klog v0.0.0-20190418165334-9cbb78b20423
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.17.1
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.17.1
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190816220812-743ec37842bf
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.17.1
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.17.1
	k8s.io/kubectl => k8s.io/kubectl v0.17.1
	k8s.io/kubelet => k8s.io/kubelet v0.17.1
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.17.1
	k8s.io/metrics => k8s.io/metrics v0.17.1
	k8s.io/node-api => k8s.io/node-api v0.17.1
	k8s.io/repo-infra => k8s.io/repo-infra v0.0.0-20181204233714-00fe14e3d1a3
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.17.1
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.17.1
	k8s.io/sample-controller => k8s.io/sample-controller v0.17.1
	k8s.io/utils => k8s.io/utils v0.0.0-20190801114015-581e00157fb1
)
