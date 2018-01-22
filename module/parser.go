package module

import (
	"strings"

	"github.com/tidwall/gjson"
	yaml2json "go.pedge.io/pkg/yaml"
)

// getBlobs extracts nested json structs from an array
// by their key
// e.g.in the yaml
//- example_rule:
//    timeout: 15s
//    match:
//      path: /foo
//    upstream:
//      name: foo_service
//      address: 127.0.0.1
//      port: 9090
//- example_rule:
//    timeout: 5s
//    match:
//      path: /bar
//    upstream:
//      name: bar_service
//      address: 127.0.0.1
//      port: 9091
// getBlobs (<yml>, "example_rule") will return
//[
//  {
//    "match": {
//      "path": "/foo"
//    },
//    "timeout": "15s",
//    "upstream": {
//      "address": "127.0.0.1",
//      "name": "foo_service",
//      "port": 9090
//    }
//  },
//  {
//    "match": {
//      "path": "/bar"
//    },
//    "timeout": "5s",
//    "upstream": {
//      "address": "127.0.0.1",
//      "name": "bar_service",
//      "port": 9091
//    }
//  }
//]

func getBlobsFromYml(yml []byte, key string) ([]byte, error) {
	jsn, err := yaml2json.ToJSON(yml, yaml2json.ToJSONOptions{})
	if err != nil {
		return nil, err
	}
	return getBlobsFromJson(jsn, key), nil
}

// expects array
func getBlobsFromJson(jsn []byte, key string) []byte {
	elements := gjson.Parse(string(jsn)).Array()
	var blobs []string
	for _, element := range elements {
		blob := element.Get(key).Raw
		if len(blob) > 0 {
			blobs = append(blobs, blob)
		}
	}
	return []byte("[" + strings.Join(blobs, ",") + "]")
}
