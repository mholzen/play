package patterns

import "github.com/fogleman/ease"

type Discretizer struct {
	start     int
	end       int
	ease      ease.Function
	intervals int
}

func NewDiscretizer(start, end int, easeFunc ease.Function, intervals int) *Discretizer {
	return &Discretizer{
		start:     start,
		end:       end,
		ease:      easeFunc,
		intervals: intervals,
	}
}

func (d *Discretizer) GetValues() []int {
	if d.intervals <= 0 {
		return []int{}
	}

	values := make([]int, d.intervals)

	for i := 0; i < d.intervals; i++ {
		t := float64(i) / float64(d.intervals-1)
		if d.intervals == 1 {
			t = 0.0
		}

		easedT := d.ease(t)
		interpolated := float64(d.start) + easedT*float64(d.end-d.start)
		values[i] = int(interpolated + 0.5) // Round to nearest integer
	}

	return values
}
