package fixture

// type Emitter[T any] struct {
// 	Value   T
// 	channel chan T
// }

// func NewEmitter[T any](value T) Emitter[T] {
// 	return Emitter[T]{Value: value}
// }

// func (e *Emitter[T]) Channel() chan T {
// 	if e.channel == nil {
// 		// log.Printf("creating channel")
// 		e.channel = make(chan T)
// 	}
// 	return e.channel
// }

// func (e *Emitter[T]) SetValue(value T) {
// 	e.Value = value
// 	e.Emit()
// }

// func (e *Emitter[T]) Emit() {
// 	select {
// 	case e.Channel() <- e.Value:
// 		// log.Printf("emitter %v emitting %v", f.Fixture, f.GetValueMap())
// 	default:
// 		// log.Printf("no receiver")
// 	}
// }
