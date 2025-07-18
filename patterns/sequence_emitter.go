package patterns

import (
	"log"

	"github.com/mholzen/play-go/controls"
)

type SequenceEmitter[T any] struct {
	sequence               *controls.Sequence[T]
	clock                  *controls.Clock
	trigger                *controls.Trigger
	*controls.Observers[T] // TODO: apply this pattern to all emitters
}

func NewSequenceEmitter[T any](sequence *controls.Sequence[T], clock *controls.Clock, trigger controls.TriggerFunc) *SequenceEmitter[T] {
	emitter := &SequenceEmitter[T]{
		sequence:  sequence,
		clock:     clock,
		Observers: controls.NewObservable[T](),
	}

	emitter.SetTriggerFunc(trigger)

	return emitter
}

func (se *SequenceEmitter[T]) SetSequence(sequence *controls.Sequence[T]) {
	se.sequence = sequence
	se.Reset()
}

func (se *SequenceEmitter[T]) SetTriggerFunc(triggerFunc controls.TriggerFunc) {
	if se.trigger != nil {
		se.clock.Cancel(se.trigger)
	}
	se.trigger = se.clock.On(triggerFunc, se.onTrigger)
	log.Printf("set trigger %p", se.trigger)
}

func (se *SequenceEmitter[T]) onTrigger() {
	value := se.sequence.Value()
	se.sequence.Inc()
	se.Notify(value)
}

func (se *SequenceEmitter[T]) Reset() {
	se.clock.Reset()
	se.sequence.Reset()
}
