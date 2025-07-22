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
