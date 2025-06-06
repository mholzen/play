package controls

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type List struct {
	items []Item `json:"-"`
}

func NewList(count int) *List {
	items := make([]Item, count)
	return &List{items: items}
}

func (l *List) SetItem(index int, item Item) {
	if index < 0 {
		panic(fmt.Sprintf("index %d is negative", index))
	}
	if index >= len(l.items) {
		newItems := make([]Item, index+1)
		copy(newItems, l.items)
		l.items = newItems
	}
	l.items[index] = item
}

func (l *List) GetItem(index string) (Item, error) {
	if index == "" {
		return nil, fmt.Errorf("index is empty") // WARN: could not be considered an error
	}

	i, err := strconv.Atoi(index)
	if err != nil {
		return nil, err
	}

	return l.items[i], nil
}

func (l *List) Items() map[string]Item {
	items := make(map[string]Item)
	for i, item := range l.items {
		items[strconv.Itoa(i)] = item
	}
	return items
}

func (l *List) Keys() []string {
	var keys []string
	for i := range l.items {
		keys = append(keys, strconv.Itoa(i))
	}
	return keys
}

func (l *List) Map() map[string]any {
	m := make(map[string]any)
	for i, item := range l.items {
		m[strconv.Itoa(i)] = item
	}
	return m
}

func (l *List) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.items)
}
