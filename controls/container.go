package controls

import "fmt"

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

func ContainerMarshalJSON(c Container) ([]byte, error) {
	return containerMarshalJSON(c, ContainerKeys(c))
}

type ItemMap map[string]Item

func (m ItemMap) GetItem(name string) (Item, error) {
	if item, ok := m[name]; ok {
		return item, nil
	}
	return nil, fmt.Errorf("item not found: '%s'", name)
}

func (m ItemMap) Items() map[string]Item {
	return m
}

func NewItemMap() ItemMap {
	return ItemMap{}
}
