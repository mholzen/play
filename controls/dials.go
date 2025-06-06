package controls

import "fmt"

type Dial[T any] interface {
	Get() T
	Set(T)
	Min() T
	Max() T
}

type DialMap[T any] map[string]Dial[T]

func NewDialMap[T any]() DialMap[T] {
	return make(DialMap[T])
}

func (d DialMap[T]) Add(name string, dial Dial[T]) {
	d[name] = dial
}

func (d DialMap[T]) Items() map[string]Item {
	items := make(map[string]Item)
	for name, dial := range d {
		items[name] = dial
	}
	return items
}

func (d DialMap[T]) GetItem(name string) (Item, error) {
	res, ok := d[name]
	if !ok {
		return nil, fmt.Errorf("dial %s not found", name)
	}
	return res, nil
}

func NewDialMapFromChannels[T any](channels []string, constructor DialConstructor[T]) DialMap[T] {
	dialMap := NewDialMap[T]()
	for _, channel := range channels {
		dialMap.Add(channel, constructor())
	}
	return dialMap
}

type DialConstructor[T any] func() Dial[T]
