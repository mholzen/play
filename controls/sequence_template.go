package controls

type SequenceT[T any] struct {
	values  []T
	Counter *Counter
}

func NewSequenceT[T any](values []T) *SequenceT[T] {
	return &SequenceT[T]{
		values:  values,
		Counter: NewCounter(len(values)),
	}
}

func (c *SequenceT[T]) Inc() {
	c.Counter.Inc()
}

func (c *SequenceT[T]) IncValues() (T, T) {
	current := c.Values()
	c.Counter.Inc()
	next := c.Values()
	return current, next
}

func (c *SequenceT[T]) Values() T {
	return c.values[c.Counter.Value()]
}

func (c *SequenceT[T]) Reset() {
	c.Counter.Reset()
}
