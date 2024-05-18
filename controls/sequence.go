package controls

type Sequence struct {
	values  []ValueMap
	Counter *Counter
}

func NewSequence(values []ValueMap) *Sequence {
	return &Sequence{
		values:  values,
		Counter: NewCounter(len(values)),
	}
}

func (c *Sequence) Inc() {
	c.Counter.Inc()
}

func (c *Sequence) IncValues() (ValueMap, ValueMap) {
	current := c.Values()
	c.Counter.Inc()
	next := c.Values()
	return current, next
}

func (c *Sequence) Values() ValueMap {
	// log.Printf("Sequence value: %d", c.Counter.Value())
	return c.values[c.Counter.Value()]
}

func (c *Sequence) Reset() {
	c.Counter.Reset()
}
