package patterns

import (
	"github.com/mholzen/play-go/controls"
)

type ChangeEmitter[T comparable] struct {
	source    controls.Observable[T]
	lastValue *T
	hasValue  bool
	*controls.Observers[T]
	observerCh chan T
}

func NewChangeEmitter[T comparable](source controls.Observable[T]) *ChangeEmitter[T] {
	emitter := &ChangeEmitter[T]{
		source:     source,
		lastValue:  nil,
		hasValue:   false,
		Observers:  controls.NewObservable[T](),
		observerCh: make(chan T),
	}

	source.AddObserver(emitter.observerCh)

	go emitter.listen()

	return emitter
}

func (ce *ChangeEmitter[T]) listen() {
	for value := range ce.observerCh {
		if !ce.hasValue || *ce.lastValue != value {
			ce.lastValue = &value
			ce.hasValue = true
			ce.Notify(value)
		}
	}
}

func (ce *ChangeEmitter[T]) Reset() {
	ce.lastValue = nil
	ce.hasValue = false
}

func (ce *ChangeEmitter[T]) Close() {
	ce.source.RemoveObserver(ce.observerCh)
}
