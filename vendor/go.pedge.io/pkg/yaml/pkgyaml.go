/*
Package pkgyaml provides functionality for YAML.

Originally inspired by https://github.com/bronze1man/yaml2json.
*/
package pkgyaml

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v2"
)

const (
	defaultIndent = "\t"
)

// ToJSONOptions are the options to pass to ToJSON.
type ToJSONOptions struct {
	// Pretty says to output the JSON with json.MarshalIndent.
	Pretty bool
	// Indent is the string to use for indenting. This only applies
	// if Pretty is set. The default is "\t".
	Indent string
}

// ToJSON transforms an YAML input and transforms it to JSON.
//
// Originally inspired by https://github.com/bronze1man/yaml2json.
func ToJSON(p []byte, opts ToJSONOptions) ([]byte, error) {
	var yamlData interface{}
	if err := yaml.Unmarshal(p, &yamlData); err != nil {
		return nil, err
	}
	jsonData, err := toJSON(yamlData)
	if err != nil {
		return nil, err
	}
	if opts.Pretty {
		if opts.Indent != "" {
			return json.MarshalIndent(jsonData, "", opts.Indent)
		}
		return json.MarshalIndent(jsonData, "", defaultIndent)
	}
	return json.Marshal(jsonData)
}

// ParseYAMLOrJSON unmarshals the given file at filePath, switching based on the file extension.
//
// This uses the json tags on any fields, as json.Unmarshal is the unmarshalling function used.
func ParseYAMLOrJSON(filePath string, v interface{}) (retErr error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	switch ext := filepath.Ext(filePath); ext {
	case ".yml", ".yaml":
		data, err = ToJSON(data, ToJSONOptions{})
		if err != nil {
			return err
		}
	case ".json":
	default:
		return fmt.Errorf("pkgyaml: %s is not a valid extension yml, yaml, or json", ext)
	}
	return json.Unmarshal(data, v)
}

func toJSON(yamlData interface{}) (interface{}, error) {
	switch yamlData.(type) {
	case map[interface{}]interface{}:
		jsonData := make(map[string]interface{})
		for key, yamlValue := range yamlData.(map[interface{}]interface{}) {
			jsonValue, err := toJSON(yamlValue)
			if err != nil {
				return nil, err
			}
			switch key.(type) {
			case string:
				jsonData[key.(string)] = jsonValue
			case int:
				jsonData[strconv.Itoa(key.(int))] = jsonValue
			default:
				return nil, fmt.Errorf("pkgyaml: unexpected key type %T for %v", key, key)
			}
		}
		return jsonData, nil
	case []interface{}:
		yamlDataSlice := yamlData.([]interface{})
		jsonDataSlice := make([]interface{}, len(yamlDataSlice))
		for i, yamlValue := range yamlDataSlice {
			jsonValue, err := toJSON(yamlValue)
			if err != nil {
				return nil, err
			}
			jsonDataSlice[i] = jsonValue
		}
		return jsonDataSlice, nil
	default:
		return yamlData, nil
	}
}
