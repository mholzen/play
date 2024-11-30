package controls

type EmitterI interface {
	SetValue(string)
	GetValue() string
	// Emit()
}

type Emitter[T any] interface {
	GetValue() T
	// SetValue(T)
	// Emit()
	Channel() <-chan T
}

type EmitterMap[T any] map[string]Emitter[T]

type Emitters[T any] []Emitter[T]

func NewEmitters[T any]() Emitters[T] {
	return make(Emitters[T], 0)
}

func (e Emitters[T]) GetString() string {
	return "<multiple items>"
}

type Receiver interface {
}
