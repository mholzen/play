package patterns

import (
	"fmt"

	"github.com/mholzen/play/controls"
)

type SequenceEmitterDial[T any] struct {
	SequenceEmitter[T]
	Duration *controls.ObservableRatioDial `json:"Duration"`
}

func (se *SequenceEmitterDial[T]) TriggerFunc() controls.TriggerFunc {
	return controls.TriggerOnPhraseRatio(se.Duration.Get().Numerator, se.Duration.Get().Denominator)
}

func NewSequenceEmitterDial[T any](sequence *controls.Sequence[T], clock *controls.Clock) *SequenceEmitterDial[T] {
	ratioDial := controls.NewObservableRatioDial()
	triggerFunc := controls.TriggerOnPhraseRatio(ratioDial.Get().Numerator, ratioDial.Get().Denominator)
	sequenceEmitter := NewSequenceEmitter(sequence, clock, triggerFunc)

	emitter := &SequenceEmitterDial[T]{
		SequenceEmitter: *sequenceEmitter,
		Duration:        ratioDial,
	}

	controls.OnChange(ratioDial, emitter.onRatioDialChange)
	return emitter
}

func (se *SequenceEmitterDial[T]) onRatioDialChange(ratio controls.Ratio) {
	triggerFunc := controls.TriggerOnPhraseRatio(ratio.Numerator, ratio.Denominator)
	se.SetTriggerFunc(triggerFunc)
	se.Reset()
}

func (se *SequenceEmitterDial[T]) Items() map[string]controls.Item {
	return map[string]controls.Item{
		"Duration": se.Duration,
	}
}

// TODO: could be drier
func (se *SequenceEmitterDial[T]) GetItem(name string) (controls.Item, error) {
	items := se.Items()
	if item, ok := items[name]; ok {
		return item, nil
	}
	return nil, fmt.Errorf("item not found")
}
