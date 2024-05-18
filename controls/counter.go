package controls

import (
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
)

// Counter struct to keep track of ticks, beats, and bars
type Counter struct {
	value  int
	Period int
	mu     sync.Mutex // Ensure thread-safe access
	Event  chan int
}

func NewCounter(period int) *Counter {
	return &Counter{
		Period: period,
		Event:  make(chan int),
	}
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = (c.value + 1) % c.Period

	// non-blocking send
	select {
	case c.Event <- c.value:
	default:
	}
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = 0
	c.Event <- c.value
}

func (c *Counter) On(values []int, callback func(int)) {
	valuesSet := mapset.NewSet[int]()
	valuesSet.Append(values...)

	go func() {
		for event := range c.Event {
			if valuesSet.Contains(event) {
				callback(event)
			}
		}
	}()
}
