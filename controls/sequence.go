package controls

type Sequence[T any] struct {
	values  []T
	Counter *Counter
}

func NewSequence[T any](values []T) *Sequence[T] {
	return &Sequence[T]{
		values:  values,
		Counter: NewCounter(len(values)),
	}
}

func (c *Sequence[T]) Inc() {
	c.Counter.Inc()
}

func (c *Sequence[T]) IncValues() (T, T) {
	current := c.Value()
	c.Counter.Inc()
	next := c.Value()
	return current, next
}

func (c *Sequence[T]) Value() T {
	return c.values[c.Counter.Value()]
}

func (c *Sequence[T]) Reset() {
	c.Counter.Reset()
}
