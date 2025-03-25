package envoy

import (
	"encoding/json"
	"fmt"
	gatewaykube "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"path/filepath"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/yaml"
	"strings"
)

type GatewayAPIOutput struct {
	OutputDir          string
	FolderPerNamespace bool
	HTTPRoutes         []*gwv1.HTTPRoute
	ListenerSets       []*ListenerSet
	RouteOptions       []*gatewaykube.RouteOption
	VirtualHostOptions []*gatewaykube.VirtualHostOption
	Upstreams          []*glookube.Upstream
	Gateway            *gwv1.Gateway
}

func (g *GatewayAPIOutput) Write() error {

	err := os.MkdirAll(g.OutputDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Gateway
	o, err := runtime.Encode(codecs.LegacyCodec(corev1.SchemeGroupVersion, gwv1.SchemeGroupVersion, gatewaykube.SchemeGroupVersion, glookube.SchemeGroupVersion), g.Gateway)
	if err != nil {
		return err
	}

	gwYaml, err2 := cleanYAMLData(o, err)
	if err2 != nil {
		return err2
	}
	folder := g.OutputDir
	if g.FolderPerNamespace {
		folder = filepath.Join(g.OutputDir, g.Gateway.Namespace)
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(fmt.Sprintf("%s/gateway-%s.yaml", folder, g.Gateway.Name), gwYaml, 0644)
	if err != nil {
		panic(err)
	}

	// Write Routes
	for _, httpRoute := range g.HTTPRoutes {
		if g.FolderPerNamespace {
			folder = filepath.Join(g.OutputDir, httpRoute.Namespace)
			err := os.MkdirAll(folder, os.ModePerm)
			if err != nil {
				return err
			}
		}
		o, err := runtime.Encode(codecs.LegacyCodec(corev1.SchemeGroupVersion, gwv1.SchemeGroupVersion, gatewaykube.SchemeGroupVersion, glookube.SchemeGroupVersion), httpRoute)
		if err != nil {
			return err
		}
		routeYaml, err2 := cleanYAMLData(o, err)
		if err2 != nil {
			return err2
		}
		err = os.WriteFile(fmt.Sprintf("%s/httproute-%s.yaml", folder, httpRoute.Name), routeYaml, 0644)
		if err != nil {
			panic(err)
		}
	}
	for _, listenerSet := range g.ListenerSets {
		if g.FolderPerNamespace {
			folder = filepath.Join(g.OutputDir, listenerSet.Namespace)
			err := os.MkdirAll(folder, os.ModePerm)
			if err != nil {
				return err
			}
		}
		listenerSetYaml, err := yaml.Marshal(listenerSet)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(fmt.Sprintf("%s/listenerset-%s.yaml", folder, listenerSet.Name), removeNullYamlFields(listenerSetYaml), 0644)
		if err != nil {
			panic(err)
		}
	}
	for _, routeOption := range g.RouteOptions {
		if g.FolderPerNamespace {
			folder = filepath.Join(g.OutputDir, routeOption.Namespace)
			err := os.MkdirAll(folder, os.ModePerm)
			if err != nil {
				return err
			}
		}
		routeOptionYaml, err := yaml.Marshal(routeOption)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(fmt.Sprintf("%s/routeoption-%s.yaml", folder, routeOption.Name), removeNullYamlFields(routeOptionYaml), 0644)
		if err != nil {
			panic(err)
		}
	}
	for _, upstream := range g.Upstreams {
		if g.FolderPerNamespace {
			folder = filepath.Join(g.OutputDir, upstream.Namespace)
			err := os.MkdirAll(folder, os.ModePerm)
			if err != nil {
				return err
			}
		}
		upstreamYaml, err := yaml.Marshal(upstream)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(fmt.Sprintf("%s/upstream-%s.yaml", folder, upstream.Name), removeNullYamlFields(upstreamYaml), 0644)
		if err != nil {
			panic(err)
		}
	}
	for _, virtualHostOptions := range g.VirtualHostOptions {
		if g.FolderPerNamespace {
			folder = filepath.Join(g.OutputDir, virtualHostOptions.Namespace)
			err := os.MkdirAll(folder, os.ModePerm)
			if err != nil {
				return err
			}
		}
		virtualHostOptionYaml, err := yaml.Marshal(virtualHostOptions)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(fmt.Sprintf("%s/virtualhostoption-%s.yaml", folder, virtualHostOptions.Name), removeNullYamlFields(virtualHostOptionYaml), 0644)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func cleanYAMLData(o []byte, err error) ([]byte, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(o, &data); err != nil {
		return nil, err
	}
	removeNullFields(data)

	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	return removeNullYamlFields(yamlData), nil
}

func removeNullYamlFields(yamlData []byte) []byte {
	stringData := strings.ReplaceAll(string(yamlData), "  creationTimestamp: null\n", "")
	stringData = strings.ReplaceAll(stringData, "status:\n", "")
	stringData = strings.ReplaceAll(stringData, "parents: null\n", "")
	stringData = strings.ReplaceAll(stringData, "status: {}\n", "")
	stringData = strings.ReplaceAll(stringData, "\n\n\n", "\n")
	stringData = strings.ReplaceAll(stringData, "\n\n", "\n")
	stringData = strings.ReplaceAll(stringData, "spec: {}\n", "")
	return []byte(stringData)
}

func (g *GatewayAPIOutput) HasItems() bool {

	if len(g.HTTPRoutes) > 0 {
		return true
	}
	if len(g.RouteOptions) > 0 {
		return true
	}
	if len(g.VirtualHostOptions) > 0 {
		return true
	}
	if len(g.Upstreams) > 0 {
		return true
	}
	// if there are only yaml objects then skip because we didnt change anything in the file

	return false
}

// removeNullFields recursively removes fields with null values from a map
func removeNullFields(m map[string]interface{}) {
	for k, v := range m {
		if v == nil {
			delete(m, k)
		} else if nestedMap, ok := v.(map[string]interface{}); ok {
			removeNullFields(nestedMap)
		} else if nestedSlice, ok := v.([]interface{}); ok {
			for _, item := range nestedSlice {
				if nestedItemMap, ok := item.(map[string]interface{}); ok {
					removeNullFields(nestedItemMap)
				}
			}
		}
	}
}

func (g *GatewayAPIOutput) ToString() (string, error) {
	output := ""

	o, err := runtime.Encode(codecs.LegacyCodec(corev1.SchemeGroupVersion, gwv1.SchemeGroupVersion, gatewaykube.SchemeGroupVersion, glookube.SchemeGroupVersion), g.Gateway)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(o, &data); err != nil {
		return "", err
	}
	removeNullFields(data)

	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}
	output += "\n---\n" + string(yamlData)

	for _, obj := range g.Upstreams {
		o, err := runtime.Encode(codecs.LegacyCodec(corev1.SchemeGroupVersion, gwv1.SchemeGroupVersion, gatewaykube.SchemeGroupVersion, glookube.SchemeGroupVersion), obj)
		if err != nil {
			return "", err
		}

		var data map[string]interface{}
		if err := json.Unmarshal(o, &data); err != nil {
			return "", err
		}
		removeNullFields(data)

		yamlData, err := yaml.Marshal(data)
		if err != nil {
			return "", err
		}
		output += "\n---\n" + string(yamlData)
	}

	for _, obj := range g.HTTPRoutes {
		o, err := runtime.Encode(codecs.LegacyCodec(corev1.SchemeGroupVersion, gwv1.SchemeGroupVersion, gatewaykube.SchemeGroupVersion, glookube.SchemeGroupVersion), obj)
		if err != nil {
			return "", err
		}

		var data map[string]interface{}
		if err := json.Unmarshal(o, &data); err != nil {
			return "", err
		}
		removeNullFields(data)

		yamlData, err := yaml.Marshal(data)
		if err != nil {
			return "", err
		}
		output += "\n---\n" + string(yamlData)
	}
	for _, op := range g.RouteOptions {
		marshaller := YamlMarshaller{}
		yaml, err := marshaller.ToYaml(&op)
		if err != nil {
			return "", err
		}

		output += "\n---\n" + string(yaml)
	}

	for _, op := range g.VirtualHostOptions {
		marshaller := YamlMarshaller{}
		yaml, err := marshaller.ToYaml(&op)
		if err != nil {
			return "", err
		}

		output += "\n---\n" + string(yaml)
	}

	// need to remove a few values
	//  creationTimestamp: null
	// status: {}
	// status:
	// parents: null
	output = strings.ReplaceAll(output, "  creationTimestamp: null\n", "")
	output = strings.ReplaceAll(output, "status:\n", "")
	output = strings.ReplaceAll(output, "parents: null\n", "")
	output = strings.ReplaceAll(output, "status: {}\n", "")
	output = strings.ReplaceAll(output, "\n\n\n", "\n")
	output = strings.ReplaceAll(output, "\n\n", "\n")
	output = strings.ReplaceAll(output, "spec: {}\n", "")

	// TODO remove leading and trailing ---
	// log.Printf("%s", output)
	return output, nil
}
