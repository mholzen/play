package patterns

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fogleman/ease"
	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
)

type RainbowControls struct {
	// controls are more than just the configurable parameters.  it also contains the dials, knobs and buttons used to select the parameters.
	Clock   *controls.Clock               `json:"-"`
	Speed   *controls.ObservableRatioDial `json:"speed"` // TODO: doesn't have to be observable
	Chase   *controls.FloatDial           `json:"chase"`
	Reverse *controls.Toggle              `json:"reverse"`
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
		// log.Printf("rainbow speed: %v", c.Speed.Value)

		ratio := c.Speed.Get().ToFloat()
		log.Printf("rainbow ratio: %v", ratio)

		duration := time.Duration(float64(c.Clock.PhrasePeriod().Nanoseconds()) * ratio)
		// log.Printf("rainbow transition duration: %v", duration)

		start, end := seq.IncValues()
		// log.Printf("transition %v to %v (duration: %v)\n", start, end, duration)

		max := len(fixtures.GetFixtureList())
		contexts := make([]context.Context, max)
		cancels := make([]context.CancelFunc, max)

		for i, f := range fixtures.GetFixtureList() {

			f.SetChannelValue("dimmer", 255)
			f.SetChannelValue("tilt", 127)

			j := i
			if c.Reverse.GetValue() {
				j = max - i - 1
			}

			if cancels[j] != nil {
				cancels[j]()
			}
			contexts[j], cancels[j] = context.WithCancel(context.Background())
			action := Transition(f, start, end, duration, ease.InOutSine, fixture.REFRESH, contexts[j])

			chaseDelay := time.Duration(float64(c.Clock.BeatPeriod().Nanoseconds()) * c.Chase.Value * float64(i))

			log.Printf("rainbow chaseDelay: %v", chaseDelay)
			go Delay(chaseDelay, action, contexts[j])
		}
	}

	t := c.Clock.On(controls.TriggerOnBar(1), transition)
	return controls.Triggers{
		*t,
	}
}

func (c RainbowControls) Items() map[string]controls.Item {
	return map[string]controls.Item{
		"speed":   c.Speed,
		"chase":   c.Chase,
		"reverse": c.Reverse,
	}
}

func (c RainbowControls) GetItem(name string) (controls.Item, error) {
	items := c.Items()
	if item, ok := items[name]; ok {
		return item, nil
	}
	return nil, fmt.Errorf("item not found")
}

func NewRainbowControls(clock *controls.Clock) RainbowControls {
	speed := controls.NewObservableRatioDial()
	chase := controls.ObservableFloatDial{ // TODO: should be a ratio of period (from 16x to 1/16x)
		FloatDial: controls.FloatDial{
			Value: 1,
			Min:   0,
			Max:   10,
		},
	}

	reverse := controls.ObservableToggle{
		Toggle: *controls.NewToggle(),
	}

	return RainbowControls{
		Clock:   clock,
		Speed:   speed,
		Chase:   &chase.FloatDial,
		Reverse: &reverse.Toggle,
	}
}
