package helpers

import "encoding/json"

func ToMap(spec interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(spec)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	return m, json.Unmarshal(data, &m)
}
