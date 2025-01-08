package controls

import "encoding/json"

type OrderedContainer interface {
	Container
	Keys() []string
}

func OrderedContainerMarshalJSON(c OrderedContainer) ([]byte, error) {
	return containerMarshalJSON(c, c.Keys())
}

func containerMarshalJSON(c Container, keys []string) ([]byte, error) {
	type Item struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	}
	res := make([]Item, 0)
	for _, key := range keys {
		item := Item{Name: key}
		value, err := c.GetItem(key)
		if err == nil {
			item.Value = value
		}
		res = append(res, item)
	}
	return json.Marshal(res)
}
