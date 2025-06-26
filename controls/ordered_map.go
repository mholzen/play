package controls

import (
	"encoding/json"
	"fmt"
)

type OrderedMap struct {
	items map[string]Item `json:"-"`
	keys  []string        `json:"-"`
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		items: make(map[string]Item),
		keys:  make([]string, 0),
	}
}

func (om *OrderedMap) SetItem(key string, item Item) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	if _, exists := om.items[key]; exists {
		return fmt.Errorf("key %q already exists", key)
	}

	om.items[key] = item
	om.keys = append(om.keys, key)
	return nil
}

func (om *OrderedMap) GetItem(key string) (Item, error) {
	if key == "" {
		return nil, fmt.Errorf("key is empty")
	}

	item, exists := om.items[key]
	if !exists {
		return nil, fmt.Errorf("key %q not found", key)
	}

	return item, nil
}

func (om *OrderedMap) Items() map[string]Item {
	items := make(map[string]Item)
	for _, key := range om.keys {
		items[key] = om.items[key]
	}
	return items
}

func (om *OrderedMap) Keys() []string {
	keys := make([]string, len(om.keys))
	copy(keys, om.keys)
	return keys
}

func (om *OrderedMap) Map() map[string]any {
	m := make(map[string]any)
	for _, key := range om.keys {
		m[key] = om.items[key]
	}
	return m
}

func (om *OrderedMap) MarshalJSON() ([]byte, error) {
	if len(om.keys) == 0 {
		return []byte("{}"), nil
	}

	var result []byte
	result = append(result, '{')

	for i, key := range om.keys {
		if i > 0 {
			result = append(result, ',')
		}

		keyBytes, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		result = append(result, keyBytes...)
		result = append(result, ':')

		valueBytes, err := json.Marshal(om.items[key])
		if err != nil {
			return nil, err
		}
		result = append(result, valueBytes...)
	}

	result = append(result, '}')
	return result, nil
}
