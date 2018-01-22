package module

import (
	"strings"

	"github.com/tidwall/gjson"
	yaml2json "go.pedge.io/pkg/yaml"
)

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
