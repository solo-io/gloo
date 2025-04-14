package snapshot

import (
	"encoding/json"
	"fmt"
	"github.com/solo-io/gloo/pkg/schemes"
	gatewaykube "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1/kube/apis/enterprise.gloo.solo.io/v1"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"log"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
)

func FromYamlFiles(files []string) (*Instance, error) {
	instance := &Instance{}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		if err := instance.translateFileToEdgeInput(string(data), file); err != nil {
			return nil, err
		}

	}

	return instance, nil
}
func FromGlooSnapshot(snapshotFile string) (*Instance, error) {
	instance := &Instance{}
	data, err := os.ReadFile(snapshotFile)
	if err != nil {
		return nil, err
	}
	if err := instance.translateSnapshotToEdgeInput(string(data), snapshotFile); err != nil {
		return nil, err
	}

	return instance, nil
}

type snapshotResponseData struct {
	Data []map[string]interface{} `json:"data"`
}

func (i *Instance) translateSnapshotToEdgeInput(jsonData string, fileName string) error {
	var snapshot snapshotResponseData
	if err := json.Unmarshal([]byte(jsonData), &snapshot); err != nil {
		return err
	}
	for _, objMap := range snapshot.Data {
		// Convert to unstructured
		obj := unstructured.Unstructured{
			Object: objMap,
		}
		if err := i.parseObjects(obj, fileName); err != nil {
			return err
		}
	}

	return nil
}

func (i *Instance) parseObjects(item unstructured.Unstructured, fileName string) error {
	runtimeScheme := runtime.NewScheme()
	// Gloo Edge APIs
	if err := schemes.SchemeBuilder.AddToScheme(runtimeScheme); err != nil {
		log.Fatal(err)
	}

	resourceYaml, err := yaml.Marshal(item)
	if err != nil {
		return err
	}
	gvk := item.GroupVersionKind()
	obj, err := runtimeScheme.New(gvk)
	if runtime.IsNotRegisteredError(err) {
		// we just want to add the yaml and move on
		if strings.Contains(gvk.Group, "solo.io") {
			i.parseErrors = append(i.ParseErrors(), fmt.Errorf("solo resource [%s] %s/%s ignored in conversion", gvk.String(), item.GetNamespace(), item.GetName()))
		}

		i.AddYamlObject(&YAMLWrapper{Yaml: string(resourceYaml), fileOrigin: fileName})
		return nil
	} else if err != nil {
		return err
	}

	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
		return fmt.Errorf("error converting unstructured to typed: %v", err)
	}

	i.addObjectToGatewayAPIOutput(obj, fileName, string(resourceYaml))
	return nil
}

func (i *Instance) addObjectToGatewayAPIOutput(obj runtime.Object, fileName string, resourceYaml string) {
	switch o := obj.(type) {
	case *glookube.Settings:
		i.AddSettings(&SettingsWrapper{Settings: o, fileOrigin: fileName})
	case *extauthv1.AuthConfig:
		i.AddAuthConfig(&AuthConfigWrapper{AuthConfig: o, fileOrigin: fileName})
	case *glookube.Upstream:
		i.AddUpstream(&UpstreamWrapper{Upstream: o, fileOrigin: fileName})
	case *gatewaykube.RouteTable:
		i.AddRouteTable(&RouteTableWrapper{RouteTable: o, fileOrigin: fileName})
	case *gatewaykube.VirtualService:
		i.AddVirtualService(&VirtualServiceWrapper{VirtualService: o, fileOrigin: fileName})
	case *gatewaykube.RouteOption:
		i.AddRouteOption(&RouteOptionWrapper{RouteOption: o, fileOrigin: fileName})
	case *gatewaykube.VirtualHostOption:
		i.AddVirtualHostOption(&VirtualHostOptionWrapper{VirtualHostOption: o, fileOrigin: fileName})
	case *gatewaykube.Gateway:
		i.AddGlooGateway(&GlooGatewayWrapper{Gateway: o, fileOrigin: fileName})
	case *gatewaykube.HttpListenerOption:
		i.AddHTTPListenerOption(&HTTPListenerOptionWrapper{HttpListenerOption: o, fileOrigin: fileName})
	case *gatewaykube.ListenerOption:
		i.AddListenerOption(&ListenerOptionWrapper{ListenerOption: o, fileOrigin: fileName})
	default:
		// if we don't know what type it is we just add it to the yaml list
		i.AddYamlObject(&YAMLWrapper{Yaml: resourceYaml, fileOrigin: fileName, Object: obj})
	}
}
func (i *Instance) translateFileToEdgeInput(yamlData string, fileName string) error {

	runtimeScheme := runtime.NewScheme()
	// Gloo Edge APIs
	if err := schemes.SchemeBuilder.AddToScheme(runtimeScheme); err != nil {
		log.Fatal(err)
	}
	codecs := serializer.NewCodecFactory(runtimeScheme)
	decoder := codecs.UniversalDeserializer()

	// Read the file
	for _, resourceYAML := range strings.Split(yamlData, "---") {
		if len(resourceYAML) == 0 {
			continue
		}
		// yaml to object
		obj, k, err := decoder.Decode([]byte(resourceYAML), nil, nil)
		if err != nil {
			if runtime.IsNotRegisteredError(err) {
				// we just want to add the yaml and move on
				i.AddYamlObject(&YAMLWrapper{Yaml: yamlData, fileOrigin: fileName})
				continue
			}
			i.parseErrors = append(i.ParseErrors(), fmt.Errorf("error parsing YAML object [%s]: %v", resourceYAML, err))
			continue
		}

		// a lot of times lists are missing the group so this object doesn't match
		if k.Kind == "List" {
			var list unstructured.UnstructuredList
			if err := yaml.Unmarshal([]byte(resourceYAML), &list); err != nil {
				return err
			}

			for _, item := range list.Items {
				err := i.parseObjects(item, fileName)
				if err != nil {
					return err
				}

			}
			continue
		}
		i.addObjectToGatewayAPIOutput(obj, fileName, resourceYAML)
	}
	return nil
}
