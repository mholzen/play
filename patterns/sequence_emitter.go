package patterns

import (
	"github.com/mholzen/play-go/controls"
)

type SequenceEmitter[T any] struct {
	sequence     *controls.SequenceT[T]
	clock        *controls.Clock
	tickDuration int
	tickCounter  int
	*controls.Observers[T]
}

func NewSequenceEmitter[T any](sequence *controls.SequenceT[T], clock *controls.Clock, tickDuration int) *SequenceEmitter[T] {
	emitter := &SequenceEmitter[T]{
		sequence:     sequence,
		clock:        clock,
		tickDuration: tickDuration,
		tickCounter:  0,
		Observers:    controls.NewObservable[T](),
	}

	clock.OnTickCallback(emitter.onTick)

	return emitter
}

func (se *SequenceEmitter[T]) onTick() {
	se.tickCounter++

	if se.tickCounter >= se.tickDuration {
		se.tickCounter = 0
		value := se.sequence.Values()
		se.sequence.Inc()
		se.Notify(value)
	}
}

func (se *SequenceEmitter[T]) Reset() {
	se.tickCounter = 0
	se.sequence.Reset()
}
