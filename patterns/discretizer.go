package patterns

import (
	"github.com/fogleman/ease"
	"golang.org/x/exp/constraints"
)

type Discretizer[T constraints.Integer] struct {
	start     T
	end       T
	ease      ease.Function
	intervals int
}

func NewDiscretizer[T constraints.Integer](start, end T, easeFunc ease.Function, intervals int) *Discretizer[T] {
	return &Discretizer[T]{
		start:     start,
		end:       end,
		ease:      easeFunc,
		intervals: intervals,
	}
}

func (d *Discretizer[T]) GetValues() []T {
	if d.intervals <= 0 {
		return []T{}
	}

	values := make([]T, d.intervals)

	for i := 0; i < d.intervals; i++ {
		t := float64(i) / float64(d.intervals-1)
		if d.intervals == 1 {
			t = 0.0
		}

		easedT := d.ease(t)
		interpolated := float64(d.start) + easedT*float64(d.end-d.start)
		values[i] = T(interpolated + 0.5) // Round to nearest integer
	}

	return values
}
