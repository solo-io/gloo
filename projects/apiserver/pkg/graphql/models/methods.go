package models

import (
	"github.com/solo-io/solo-kit/pkg/errors"
)

func (m InputMapStringString) Validate() error {
	counted := make(map[string]int)
	for _, val := range m.Values {
		counted[val.Key]++
		if counted[val.Key] > 1 {
			return errors.Errorf("key %v appared twice in MapStringString", val.Key)
		}
	}
	return nil
}

func (m InputMapStringString) GoType() map[string]string {
	goMap := make(map[string]string)
	for _, val := range m.Values {
		goMap[val.Key] = val.Value
	}
	return goMap
}

func NewMapStringString(m map[string]string) *MapStringString {
	if len(m) == 0 {
		return nil
	}
	var values []Value
	for k, v := range m {
		values = append(values, Value{Key: k, Value: v})
	}
	return &MapStringString{
		Values: values,
	}
}
