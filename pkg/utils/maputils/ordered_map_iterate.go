package maputils

import (
	"sort"
)

// thank you R.P.
// TODO(ilackarms): find a better way to solve generics than interface{}
// consider using solo-kit to generate
func OrderedMapIterator(m map[string]interface{}, onKey func(key string, value interface{})) {
	var list []struct {
		key   string
		value interface{}
	}
	for k, v := range m {
		list = append(list, struct {
			key   string
			value interface{}
		}{
			key:   k,
			value: v,
		})
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].key < list[j].key
	})
	for _, el := range list {
		onKey(el.key, el.value)
	}
}
