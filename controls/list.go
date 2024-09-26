package controls

import "strconv"

type List struct {
	items []Item `json:"-"`
}

func NewList(count int) *List {
	items := make([]Item, count)
	return &List{items: items}
}

func (l *List) SetItem(index int, item Item) {
	l.items[index] = item
}

func (l *List) GetItem(index string) (Item, error) {
	i, err := strconv.Atoi(index)
	if err != nil {
		return nil, err
	}

	return l.items[i], nil
}

func (l *List) Items() []string {
	var items []string
	for _, item := range l.items {
		items = append(items, item.GetString())
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