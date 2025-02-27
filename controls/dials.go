package controls

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

func NewObservableDialMap(channels []string) DialMap[byte] {
	dialMap := NewDialMap[byte]()
	for _, channel := range channels {
		dialMap.Add(channel, NewObservableNumericalDial(NewNumericDial()))
	}
	return dialMap
}

func NewObservableDialMap3(channels []string) *ObservableDialMap2 {
	res := NewObservableDialMap2()
	for _, channel := range channels {
		res.AddItem(channel, NewObservableNumericalDial(NewNumericDial()))
	}
	return res
}

func ChannelsToDialMap[T any](channels []string, constructor DialConstructor[T]) DialMap[T] {
	dialMap := NewDialMap[T]()
	for _, channel := range channels {
		dialMap.Add(channel, constructor(channel))
	}
	return dialMap
}

func ChannelsToDialMap2[T any](channels []string, constructor func() Dial[T]) DialMap[T] {
	dialMap := NewDialMap[T]()
	for _, channel := range channels {
		dialMap.Add(channel, constructor())
	}
	return dialMap
}

type DialConstructor[T any] func(channel string) Dial[T]
