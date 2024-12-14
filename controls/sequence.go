package controls

type Sequence struct {
	values  []ChannelValues
	Counter *Counter
}

func NewSequence(values []ChannelValues) *Sequence {
	return &Sequence{
		values:  values,
		Counter: NewCounter(len(values)),
	}
}

func (c *Sequence) Inc() {
	c.Counter.Inc()
}

func (c *Sequence) IncValues() (ChannelValues, ChannelValues) {
	current := c.Values()
	c.Counter.Inc()
	next := c.Values()
	return current, next
}

func (c *Sequence) Values() ChannelValues {
	// log.Printf("Sequence value: %d", c.Counter.Value())
	return c.values[c.Counter.Value()]
}

func (c *Sequence) Reset() {
	c.Counter.Reset()
}
