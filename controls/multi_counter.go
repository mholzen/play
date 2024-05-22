package controls

// Counter struct to keep track of ticks, beats, and bars
type MultiCounter []Counter

func (c *MultiCounter) Inc() {
	for i := range *c {
		(*c)[i].Inc()
		if (*c)[i].Value() != 0 {
			break
		}
	}
}

func (c *MultiCounter) Values() []int {
	values := make([]int, len(*c))
	for i := range *c {
		values[i] = (*c)[i].Value()
	}
	return values
}

func NewMultiCounter(periods ...int) MultiCounter {
	counters := make(MultiCounter, len(periods))
	for i, period := range periods {
		counters[i] = *NewCounter(period)
	}
	return counters
}

type MultiCounter2 struct {
	MultiCounter MultiCounter
	Names        []string
}

func (c *MultiCounter2) On(name string, counts []int, fn func(int)) {
	for i := range c.Names {
		if c.Names[i] == name {
			c.MultiCounter[i].On(counts, fn)
		}
	}
}

// var clock = MultiCounter2{
// 	MultiCounter: NewMultiCounter(4, 4, 24),
// 	Names:        []string{"bar", "beat", "tick"},
// }
