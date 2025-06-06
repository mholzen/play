package controls

import (
	"fmt"
)

type Container interface {
	Item
	GetItem(string) (Item, error) // TODO: should return (Item, bool)
	Items() map[string]Item
}

func ContainerKeys(c Container) []string {
	keys := make([]string, 0, len(c.Items()))
	for k := range c.Items() {
		keys = append(keys, k)
	}
	return keys
}

func ContainerFollowPath(c Container, segments []string) (Item, error) {
	for i, segment := range segments {
		// log.Printf("segment: '%s'", segment)
		item, err := c.GetItem(segment)
		if err != nil {
			return nil, err
		}
		if i == len(segments)-1 {
			return item, nil
		}
		cl, ok := item.(Container)
		if !ok {
			return nil, fmt.Errorf("item is not a Container (got %T)", item)
		}
		c = cl
	}
	return c, nil
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
