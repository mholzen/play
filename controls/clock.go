package controls

import (
	"fmt"
	"math"
	"time"
)

// Clock struct definition
type Clock struct {
	bpm            float64
	ticks          int
	beats          int
	bars           int
	counter        Counter
	beat, bar, pos int
	start          time.Time
	intervalMillis float64
	ticker         *time.Ticker
	stopChan       chan struct{}
}

func NewCustomClock(bpm float64, ticks, beats, bars int) *Clock {
	if bpm == 0 {
		bpm = 120
	}
	clock := &Clock{
		bpm:      bpm,
		ticks:    ticks,
		beats:    beats,
		bars:     bars,
		start:    time.Now(),
		stopChan: make(chan struct{}),
	}
	clock.intervalMillis = (60.0 / clock.bpm) / float64(clock.ticks) * 1000.0
	return clock
}

func NewClock(bpm float64) *Clock {
	return NewCustomClock(bpm, 24, 4, 4)
}

func (c *Clock) scheduleNext() {
	now := time.Now()
	nextInterval := time.Duration(math.Round(c.intervalMillis*1000)) * time.Microsecond
	nextTickTime := c.start.Add(nextInterval * time.Duration(c.counter.Value()+1))
	if nextTickTime.Before(now) {
		nextInterval = 0
	}
	c.ticker = time.NewTicker(nextTickTime.Sub(now))
	go func() {
		for {
			select {
			case <-c.ticker.C:
				c.tick()
			case <-c.stopChan:
				c.ticker.Stop()
				return
			}
		}
	}()
}

func (c *Clock) Start() {
	c.scheduleNext()
}

func (c *Clock) Stop() {
	close(c.stopChan)
}

func (c *Clock) Reset() {
	c.Stop() // Stop the current ticking
	c.counter.Reset()
	c.beat, c.bar, c.pos = 0, 0, 0
	c.start = time.Now()
	c.scheduleNext()
}

func (c *Clock) tick() {
	c.counter.Inc()
	c.pos = c.counter.Value()
	// Emit tick - Here you would send the tick event through a channel or call a callback function
	fmt.Printf("Tick at %d\n", c.pos)
	// Similar logic for beats, bars, and other events as in the CoffeeScript version
	c.scheduleNext()
}
