package patterns

import (
	"fmt"
	"log"

	"github.com/fogleman/ease"
	"github.com/mholzen/play-go/controls"
)

type DialPattern struct {
	Mux *controls.Mux[byte]

	Dial *controls.ObservableNumericalDial

	SequenceEmitter *SequenceEmitterDial[byte]

	EaseSelector *controls.Selector[ease.Function]
	Intervals    *controls.ObservableNumericalDial
}

func NewDialPattern(clock *controls.Clock) *DialPattern {
	mux := controls.NewMux[byte]()

	dial := controls.NewObservableNumericalDial()
	mux.Add("dial", dial)

	intervals := controls.NewObservableNumericalDial()
	intervals.Set(16)

	easeSelector := controls.NewSelector[ease.Function]()
	easeSelector.SetOptions(EaseMap)
	easeSelector.SetSelected("square")

	sequence := NewSequence(easeSelector.GetSelectedValue(), int(intervals.Get()))
	sequenceEmitter := NewSequenceEmitterDial(sequence, clock)

	mux.Add("pattern", sequenceEmitter)

	controls.OnChange(easeSelector, func(ease ease.Function) {
		sequence := NewSequence(ease, int(intervals.Get()))
		sequenceEmitter.SetSequence(sequence)
	})

	controls.OnChange(intervals, func(value byte) {
		if value <= 0 {
			log.Printf("ignoring intervals set to 0 or negative (%d)", value)
			return
		}
		sequence := NewSequence(easeSelector.GetSelectedValue(), int(value))
		sequenceEmitter.SetSequence(sequence)
	})

	return &DialPattern{
		Dial:            dial,
		EaseSelector:    easeSelector,
		Intervals:       intervals,
		SequenceEmitter: sequenceEmitter,
		Mux:             mux,
	}
}

func NewSequence(ease ease.Function, intervals int) *controls.Sequence[byte] {
	discretizer := NewDiscretizer(byte(0), byte(255), ease, intervals)
	values := discretizer.GetValues()
	sequence := controls.NewSequence(values)
	return sequence
}

func (c *DialPattern) Items() map[string]controls.Item {
	return map[string]controls.Item{
		"Mux":             c.Mux,
		"Dial":            c.Dial,
		"SequenceEmitter": c.SequenceEmitter,
		"EaseSelector":    c.EaseSelector,
		"Intervals":       c.Intervals,
	}
}

func (c *DialPattern) GetItem(name string) (controls.Item, error) {
	item, ok := c.Items()[name]
	if !ok {
		return nil, fmt.Errorf("item %s not found", name)
	}
	return item, nil
}
