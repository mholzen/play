package controls

import (
	"fmt"
	"log"
	"time"
)

// Clock struct definition
type Clock struct {
	Bpm          float64
	TicksPerBeat int
	BeatsPerBar  int
	BarPerPhrase int

	ticks    int
	start    time.Time
	ticker   *time.Ticker
	stopChan chan struct{}

	TickC chan struct{}
	BeatC chan struct{}
	BarC  chan struct{}

	tickCallbacks []func()
	Triggers      []Trigger

	NudgeDelay  time.Duration
	ResetPeriod bool
}

func (c *Clock) Ticks() int {
	return c.ticks
}
func (c *Clock) Tick() int {
	return c.ticks % c.TicksPerBeat
}
func (c *Clock) Beats() int {
	return int(c.ticks / c.TicksPerBeat)
}
func (c *Clock) Bars() int {
	return int(c.Beats() / c.BeatsPerBar)
}
func (c *Clock) Phrases() int {
	return int(c.Bars() / c.BarPerPhrase)
}

func (c *Clock) Beat() int {
	return c.Beats() % c.BeatsPerBar
}
func (c *Clock) Bar() int {
	return c.Bars() % c.BarPerPhrase
}
func (c *Clock) Phrase() int {
	return c.Phrases()
}

func (c *Clock) TickPeriod() time.Duration {
	ticksPerMin := c.Bpm * float64(c.TicksPerBeat)
	return time.Duration(60000.0/ticksPerMin) * time.Millisecond
}
func (c *Clock) BeatPeriod() time.Duration {
	return c.TickPeriod() * time.Duration(c.TicksPerBeat)
}
func (c *Clock) BarPeriod() time.Duration {
	return c.BeatPeriod() * time.Duration(c.BeatsPerBar)
}
func (c *Clock) PhrasePeriod() time.Duration {
	return c.BarPeriod() * time.Duration(c.BarPerPhrase)
}

func (c *Clock) Start() {
	c.start = time.Now()
	c.ticker = time.NewTicker(c.TickPeriod())

	c.Trigger() // Trigger 0

	go func() {
		for {
			select {
			case <-c.ticker.C:
				c.tick()

				if c.ResetPeriod {
					c.ticker = time.NewTicker(c.TickPeriod())
					c.ResetPeriod = false
				}

				if c.NudgeDelay != 0 {
					c.ticker = time.NewTicker(c.TickPeriod() + c.NudgeDelay)
					c.NudgeDelay = 0
					c.ResetPeriod = true
				}

			case <-c.stopChan:
				c.ticker.Stop()
				return
			}
		}
	}()
}

func (c *Clock) Stop() {
	close(c.stopChan)
}

func (c *Clock) Reset() {
	// c.Stop() // Stop the current ticker
	c.ticks = 0
}
func (c *Clock) SetBpm(bpm float64) {
	c.Bpm = bpm
	c.ticker = time.NewTicker(c.TickPeriod())
}

func (c *Clock) Nudge(delta time.Duration) {
	c.NudgeDelay = c.TickPeriod() + delta
	c.ResetPeriod = false
}

func (c *Clock) Pace(delta float64) {
	c.Bpm += delta
	c.ResetPeriod = false
}

func (c *Clock) String() string {
	return fmt.Sprintf("Bpm: %f Phrase: %d Bar: %d Beat: %d\n",
		c.Bpm,
		c.Phrase(),
		c.Bar(),
		c.Beat(),
	)
}

func (c *Clock) StringLong() string {
	return fmt.Sprintf("Time: %s Ticks: %d Beats: %d Bars: %d -- Tick: %d Beat: %d Bar: %d\n",
		time.Now(),
		c.ticks, c.Beats(), c.Bars(),
		c.Tick(), c.Beat(), c.Bar(),
	)
}
func (c *Clock) tick() {
	c.ticks++
	c.Trigger()
}

func (c *Clock) Trigger() {
	for _, callback := range c.tickCallbacks {
		callback()
	}
	// c.SendToChannels()
	c.CheckTriggers()
}

func (c *Clock) On(trigger TriggerFunc, callback func()) *Trigger {
	t := Trigger{trigger, true, callback}
	c.Triggers = append(c.Triggers, t)
	log.Printf("added trigger %v\n", t)
	return &t
}

func (c *Clock) CheckTriggers() {
	for _, trigger := range c.Triggers {
		if trigger.When(*c) && trigger.Enabled {
			go trigger.Do()
		}
	}
}

func (c *Clock) SendToChannels() {
	select {
	case c.TickC <- struct{}{}:
	default:
	}

	if c.Beat() == 0 {
		select {
		case c.BeatC <- struct{}{}:
		default:
		}
	}
	if c.Bar() == 0 {
		select {
		case c.BarC <- struct{}{}:
		default:
		}
	}
}

func (c *Clock) OnTickCallback(callback func()) {
	c.tickCallbacks = append(c.tickCallbacks, callback)
}

func NewCustomClock(bpm float64, ticks, beats, bars, phrases int) *Clock {
	if bpm == 0 {
		bpm = 120
	}
	clock := &Clock{
		Bpm:          bpm,
		TicksPerBeat: ticks,
		BeatsPerBar:  beats,
		BarPerPhrase: beats,

		stopChan: make(chan struct{}),

		TickC: make(chan struct{}),
		BeatC: make(chan struct{}),
		BarC:  make(chan struct{}),
	}
	return clock
}

func NewClock(bpm float64) *Clock {
	return NewCustomClock(bpm, 24, 4, 4, 4)
}
