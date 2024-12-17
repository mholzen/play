package controls

import (
	"encoding/json"
	"fmt"
)

type Container interface {
	Item
	GetItem(string) (Item, error)
	Items() map[string]Item
}

func ContainerKeys(c Container) []string {
	keys := make([]string, 0, len(c.Items()))
	for k := range c.Items() {
		keys = append(keys, k)
	}
	return keys
}

type Map struct {
	items map[string]Item
}

func (m Map) GetItem(key string) (Item, error) {
	item, ok := m.items[key]
	if !ok {
		return nil, fmt.Errorf("item not found")
	}
	return item, nil
}

func (m Map) Items() map[string]Item {
	return m.items
}

func NewMap(items ...Item) *Map {
	res := Map{
		items: make(map[string]Item),
	}
	return &res
}

func (c *Map) AddItem(name string, item Item) {
	c.items[name] = item
}

func (m Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.items)
}
