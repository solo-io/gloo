package helpers

import (
	"os/exec"
	"reflect"
	"strings"
)

func GlooDirectory() string {
	data, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(data))
}

// Merge maps, merge priority is left to right (i.e. right maps will overwrite left maps)
func MergeValueMaps(maps ...map[string]interface{}) map[string]interface{} {
	var out map[string]interface{}
	for i, m := range maps {
		if i == 0 {
			out = m
			continue
		}
		out = mergeTwoValueMaps(out, m)
	}

	return out
}

// Yanked this right out of helm libs, merge priority is left to right
func mergeTwoValueMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := CastMap(v); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := CastMap(bv); ok {
					out[k] = mergeTwoValueMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

// Made our own cast map because regular casting was yielding inconsistent results
func CastMap(value interface{}) (map[string]interface{}, bool) {
	result := map[string]interface{}{}
	rValues := reflect.ValueOf(value)
	if rValues.Kind() == reflect.Map {
		iter := rValues.MapRange()
		for iter.Next() {
			k := iter.Key().Interface()
			v := iter.Value().Interface()
			kS, isS := k.(string)
			if isS {
				result[kS] = v
			} else {
				return map[string]interface{}{}, false

			}
		}
		return result, true
	} else {
		return map[string]interface{}{}, false
	}
}
