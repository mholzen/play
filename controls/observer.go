package controls

type Observer[T any] func(newValue T)

type Observable[T any] struct {
	value     T
	observers []Observer[T]
}

func (o *Observable[T]) Set(value T) {
	o.value = value
	o.notifyObservers(value)
}

func (o *Observable[T]) Get() T {
	return o.value
}

func (o *Observable[T]) AddObserver(observer Observer[T]) {
	o.observers = append(o.observers, observer)
}

func (o *Observable[T]) notifyObservers(value T) {
	for _, observer := range o.observers {
		observer(value)
	}
}
