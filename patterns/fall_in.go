package patterns

import (
	"context"
	"time"

	"github.com/fogleman/ease"
	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
)

type FallInControls struct {
	Clock *controls.Clock `json:"-"`
}

func (c FallInControls) FallIn(fixtures *fixture.AddressableFixtures[fixture.Fixture]) controls.Triggers {

	down := c.Clock.On(controls.TriggerOnBeat(0), func() {
		// log.Printf("fall in: start down")
		fixtures.SetChannelValue("dimmer", 255)
		fixtures.SetChannelValue("speed", 240) // should be a function to compute the speed needed to get to the tilt=0 by beat 3

		fixtures.SetChannelValue("tilt", 0) // target tilt -- will take time

		dimmerTop, _ := controls.NewChannelValues("dimmer:255")
		dimmerBottom, _ := controls.NewChannelValues("dimmer:0")

		go Transition(fixtures, dimmerTop, dimmerBottom, c.Clock.BeatPeriod()*3, ease.OutSine, 10*time.Millisecond, context.Background())()
	})

	up := c.Clock.On(controls.TriggerOnBeat(3), func() {
		// log.Printf("fall in: start up")
		fixtures.SetChannelValue("dimmer", 0)
		fixtures.SetChannelValue("speed", 0)
		fixtures.SetChannelValue("tilt", 128)
	})

	return controls.Triggers{
		*down,
		*up,
	}
}
