package controls

import (
	"encoding/json"
	"strconv"
)

type Labelable interface {
	comparable
	Label() string
}

type Options[T comparable] struct {
	Value T `json:"value"`
}

type DiscreteDial[T Labelable] struct {
	Index   int          `json:"index"`
	Options []Options[T] `json:"options"`
}

func NewDiscreteDial[T Labelable](options []T) *DiscreteDial[T] {
	labeledOptions := make([]Options[T], len(options))
	for i, option := range options {
		labeledOptions[i] = Options[T]{
			Value: option,
		}
	}
	return &DiscreteDial[T]{
		Index:   0,
		Options: labeledOptions,
	}
}

func (d *DiscreteDial[T]) Get() T {
	return d.Options[d.Index].Value
}

func (d *DiscreteDial[T]) GetMark() Options[T] {
	return d.Options[d.Index]
}

func (d *DiscreteDial[T]) Set(value T) {
	for i, option := range d.Options {
		if option.Value == value {
			d.Index = i
			return
		}
	}
	panic("value not found")
}

func (d *DiscreteDial[T]) SetIndex(index int) {
	d.Index = index
}

func (d *DiscreteDial[T]) GetValueString() string {
	return d.Get().Label()
}

func (d *DiscreteDial[T]) SetValueString(index string) error {
	indexInt, err := strconv.Atoi(index)
	if err != nil {
		return err
	}
	d.SetIndex(indexInt)
	return nil
}

func (d *DiscreteDial[T]) MarshalJSON() ([]byte, error) {
	if len(d.Options) == 0 {
		return json.Marshal(d)
	}

	marks := make([]sliderMarksJSON, len(d.Options))
	for i, option := range d.Options {
		marks[i] = sliderMarksJSON{
			Value: i,
			Label: option.Value.Label(),
		}
	}

	return json.Marshal(sliderJSON[T]{
		Value: d.Index,
		Marks: marks,
		Min:   0,
		Max:   len(d.Options) - 1,
	})
}

type sliderMarksJSON struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}

type sliderJSON[T comparable] struct {
	Value int               `json:"value"`
	Marks []sliderMarksJSON `json:"marks"`
	Min   int               `json:"min"`
	Max   int               `json:"max"`
}
