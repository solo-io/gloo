package istioconverter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	meshconfig "istio.io/api/mesh/v1alpha1"
	routing "istio.io/api/routing/v1alpha1"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/serviceregistry/kube"
	"istio.io/istio/pkg/log"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func createIngressRule(host, path string, backend v1beta1.IngressBackend, tlsSecret string) *routing.IngressRule {
	rule := &routing.IngressRule{
		Destination: &routing.IstioService{
			Name: backend.ServiceName,
		},
		TlsSecret: tlsSecret,
		Match: &routing.MatchCondition{
			Request: &routing.MatchRequest{
				Headers: make(map[string]*routing.StringMatch, 2),
			},
		},
	}
	switch backend.ServicePort.Type {
	case intstr.Int:
		rule.DestinationServicePort = &routing.IngressRule_DestinationPort{
			DestinationPort: int32(backend.ServicePort.IntValue()),
		}
	case intstr.String:
		rule.DestinationServicePort = &routing.IngressRule_DestinationPortName{
			DestinationPortName: backend.ServicePort.String(),
		}
	}

	if host != "" {
		rule.Match.Request.Headers[model.HeaderAuthority] = &routing.StringMatch{
			MatchType: &routing.StringMatch_Exact{Exact: host},
		}
	}

	if path != "" {
		if strings.HasSuffix(path, ".*") {
			rule.Match.Request.Headers[model.HeaderURI] = &routing.StringMatch{
				MatchType: &routing.StringMatch_Prefix{Prefix: strings.TrimSuffix(path, ".*")},
			}
		} else {
			rule.Match.Request.Headers[model.HeaderURI] = &routing.StringMatch{
				MatchType: &routing.StringMatch_Exact{Exact: path},
			}
		}
	} else {
		rule.Match.Request.Headers[model.HeaderURI] = &routing.StringMatch{
			MatchType: &routing.StringMatch_Prefix{Prefix: "/"},
		}
	}
	return rule
}

// encodeIngressRuleName encodes an ingress rule name for a given ingress resource name,
// as well as the position of the rule and path specified within it, counting from 1.
// ruleNum == pathNum == 0 indicates the default backend specified for an ingress.
func encodeIngressRuleName(ingressName string, ruleNum, pathNum int) string {
	return fmt.Sprintf("%s-%d-%d", ingressName, ruleNum, pathNum)
}

// decodeIngressRuleName decodes an ingress rule name previously encoded with encodeIngressRuleName.
func decodeIngressRuleName(name string) (ingressName string, ruleNum, pathNum int, err error) {
	parts := strings.Split(name, "-")
	if len(parts) < 3 {
		err = fmt.Errorf("could not decode string into ingress rule name: %s", name)
		return
	}

	ingressName = strings.Join(parts[0:len(parts)-2], "-")
	ruleNum, ruleErr := strconv.Atoi(parts[len(parts)-2])
	pathNum, pathErr := strconv.Atoi(parts[len(parts)-1])

	if pathErr != nil || ruleErr != nil {
		err = multierror.Append(
			fmt.Errorf("could not decode string into ingress rule name: %s", name),
			pathErr, ruleErr)
		return
	}

	return
}

// shouldProcessIngress determines whether the given ingress resource should be processed
// by the controller, based on its ingress class annotation.
// See https://github.com/kubernetes/ingress/blob/master/examples/PREREQUISITES.md#ingress-class
func shouldProcessIngress(mesh *meshconfig.MeshConfig, ingress *v1beta1.Ingress) bool {
	class, exists := "", false
	if ingress.Annotations != nil {
		class, exists = ingress.Annotations[kube.IngressClassAnnotation]
	}

	switch mesh.IngressControllerMode {
	case meshconfig.MeshConfig_OFF:
		return false
	case meshconfig.MeshConfig_STRICT:
		return exists && class == mesh.IngressClass
	case meshconfig.MeshConfig_DEFAULT:
		return !exists || class == mesh.IngressClass
	default:
		log.Warnf("invalid ingress synchronization mode: %v", mesh.IngressControllerMode)
		return false
	}
}
