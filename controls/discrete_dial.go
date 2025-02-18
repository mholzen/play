package controls

import (
	"encoding/json"
)

type DiscreteDial[T comparable] struct {
	Index int `json:"index"`
	// Value T   `json:"value"`
	Marks []T `json:"marks"`
}

func NewDiscreteDial[T comparable](marks []T) *DiscreteDial[T] {
	return &DiscreteDial[T]{
		Index: 0,
		Marks: marks,
	}
}

type discreteDialJSON[T comparable] struct {
	Value T   `json:"value"`
	Marks []T `json:"marks"`
	Min   int `json:"min"`
	Max   int `json:"max"`
}

func (d *DiscreteDial[T]) Get() T {
	return d.Marks[d.Index]
}

func (d *DiscreteDial[T]) Set(value T) {
	// TODO: avoid failing silently if the value is not in the marks
	for i, mark := range d.Marks {
		if mark == value {
			d.Index = i
			// d.Value = value
			return
		}
	}
}

func (d *DiscreteDial[T]) MarshalJSON() ([]byte, error) {
	if len(d.Marks) == 0 {
		return json.Marshal(d)
	}

	return json.Marshal(discreteDialJSON[T]{
		Value: d.Get(),
		Marks: d.Marks,
		Min:   0,
		Max:   len(d.Marks) - 1,
	})
}
