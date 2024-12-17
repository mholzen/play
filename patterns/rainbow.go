package patterns

import (
	"log"
	"time"

	"github.com/fogleman/ease"
	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
)

type RainbowControls struct {
	Clock *controls.Clock
	Speed *controls.FloatDial
	Chase *controls.FloatDial
}

func (c RainbowControls) Rainbow(fixtures *fixture.AddressableFixtures[fixture.Fixture]) controls.Triggers {
	seq := controls.NewSequence([]controls.ChannelValues{
		controls.AllColors["red"].Values(), // TODO: if `red` doesn't exist, this should fail fast rather than return a the 0 (ie. black) color
		controls.AllColors["yellow"].Values(),
		controls.AllColors["green"].Values(),
		controls.AllColors["cyan"].Values(),
		controls.AllColors["blue"].Values(),
		controls.AllColors["purple"].Values(),
	})

	transition := func() {
		duration := time.Duration(float64(c.Clock.PhrasePeriod().Nanoseconds()) / c.Speed.Value)
		start, end := seq.IncValues()
		log.Printf("transition %v to %v (duration: %v)\n", start, end, duration)
		for i, f := range fixtures.GetFixtureList() {
			action := Transition(f, start, end, duration, ease.InOutSine, fixture.REFRESH)

			f.SetChannelValue("dimmer", 255)
			f.SetChannelValue("tilt", 127)

			d := time.Duration(i) * (duration / 2)
			go Delay(d, action)
		}
	}

	t := c.Clock.On(controls.TriggerOnBar(1), transition)
	return controls.Triggers{
		*t,
	}
}

func (c *RainbowControls) GetContainer() controls.Container {
	dialMap := controls.NewMap()
	dialMap.AddItem("speed", c.Speed)
	dialMap.AddItem("chase", c.Chase)
	return dialMap
}

type ObservableFloatDial struct {
	controls.Observers[controls.FloatDial]
	controls.FloatDial
}

func NewRainbowControls(clock *controls.Clock) RainbowControls {

	speed := ObservableFloatDial{
		FloatDial: controls.FloatDial{
			Value: 1,
			Min:   -10,
			Max:   10,
		},
		Observers: *controls.NewObservable[controls.FloatDial](),
	}
	speed.Observers.AddObserverFunc(func(value controls.FloatDial) {
		speed.FloatDial = value
	})

	chase := ObservableFloatDial{
		FloatDial: controls.FloatDial{
			Value: 0,
			Min:   0,
			Max:   10,
		},
		Observers: *controls.NewObservable[controls.FloatDial](),
	}
	chase.Observers.AddObserverFunc(func(value controls.FloatDial) {
		chase.FloatDial = value
	})

	return RainbowControls{
		Clock: clock,
		Speed: &speed.FloatDial,
		Chase: &chase.FloatDial,
	}
}
