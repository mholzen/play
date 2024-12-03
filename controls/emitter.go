package controls

type EmitterI interface {
	SetValue(string)
	GetValue() string
	// Emit()
}

type Emitter[T any] interface {
	GetValue() T
	// SetValue(T)
	Emit()
	Channel() <-chan T
}

type EmitterMap[T any] map[string]EmitterValue[T]

type Emitters[T any] []EmitterValue[T]

func NewEmitters[T any]() Emitters[T] {
	return make(Emitters[T], 0)
}

func (e Emitters[T]) GetString() string {
	return "<multiple items>"
}

type Receiver interface {
}

type EmitterValue[T any] struct {
	Value   T
	channel chan T
}

func NewEmitterValue[T any](value T) EmitterValue[T] {
	return EmitterValue[T]{Value: value}
}

func (e *EmitterValue[T]) Channel() <-chan T {
	if e.channel == nil {
		// log.Printf("creating channel")
		e.channel = make(chan T)
	}
	return e.channel
}

func (e *EmitterValue[T]) SetValue(value T) {
	e.Value = value
	e.Emit()
}

func (e *EmitterValue[T]) GetValue() T {
	return e.Value
}

func (e *EmitterValue[T]) Emit() {
	if e.channel == nil {
		e.channel = make(chan T)
	}
	select {
	case e.channel <- e.Value:
		// log.Printf("emitter %v emitting %v", f.Fixture, f.GetValueMap())
	default:
		// log.Printf("no receiver")
	}
}
