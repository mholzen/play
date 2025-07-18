package patterns

import (
	"context"
	"fmt"
	"time"

	"github.com/fogleman/ease"
	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
)

type RainbowControls struct {
	// controls are more than just the configurable parameters.  it also contains the dials, knobs and buttons used to select the parameters.
	Clock   *controls.Clock               `json:"-"`
	Cycle   *controls.ObservableRatioDial `json:"cycle"`
	Speed   *controls.ObservableRatioDial `json:"speed"` // TODO: doesn't have to be observable
	Chase   *controls.FloatDial           `json:"chase"`
	Reverse *controls.Toggle              `json:"reverse"`
}

func (c RainbowControls) Rainbow(fixtures *fixture.AddressableFixtures[fixture.Fixture]) {
	// initial setup -- could be optional

	fixtures.SetChannelValue("dimmer", 255)
	fixtures.SetChannelValue("tilt", 127)

	seq := controls.NewSequenceT([]controls.ChannelValues{
		controls.AllColors["red"].Values(), // TODO: if `red` doesn't exist, this should fail fast rather than return a the 0 (ie. black) color
		controls.AllColors["yellow"].Values(),
		controls.AllColors["green"].Values(),
		controls.AllColors["cyan"].Values(),
		controls.AllColors["blue"].Values(),
		controls.AllColors["purple"].Values(),
	})

	max := len(fixtures.GetFixtureList())
	contexts := make([]context.Context, max)
	cancels := make([]context.CancelFunc, max)

	transition := func() {
		ratio := c.Speed.Get().ToFloat()
		// log.Printf("rainbow ratio: %v", ratio)

		duration := time.Duration(float64(c.Clock.PhrasePeriod().Nanoseconds()) * ratio)
		// log.Printf("rainbow transition duration: %v", duration)

		start, end := seq.IncValues()
		// log.Printf("transition %v to %v (duration: %v)\n", start, end, duration)

		for i, f := range fixtures.GetFixtureList() {

			var fixtureIndex int
			if c.Reverse.GetValue() {
				fixtureIndex = max - i - 1
			} else {
				fixtureIndex = i
			}

			if cancelFunc := cancels[fixtureIndex]; cancelFunc != nil {
				cancelFunc()
			}

			contexts[fixtureIndex], cancels[fixtureIndex] = context.WithCancel(context.Background())
			action := Transition(f, start, end, duration, ease.InOutSine, fixture.REFRESH, contexts[fixtureIndex])

			chaseDelay := time.Duration(float64(c.Clock.BeatPeriod().Nanoseconds()) * c.Chase.Value * float64(i))

			// log.Printf("rainbow chaseDelay: %v", chaseDelay)
			go Delay(chaseDelay, action, contexts[fixtureIndex])
		}
	}

	var trigger *controls.Trigger
	setupTrigger := func(ratio controls.Ratio) {
		if trigger != nil {
			c.Clock.Cancel(trigger)
		}
		trigger = c.Clock.On(controls.TriggerOnPhraseRatio(ratio.Numerator, ratio.Denominator), transition)
	}

	setupTrigger(c.Cycle.Get())
	controls.OnChange(c.Cycle, setupTrigger)
}

func (c RainbowControls) Items() map[string]controls.Item {
	return map[string]controls.Item{
		"speed":   c.Speed,
		"chase":   c.Chase,
		"reverse": c.Reverse,
		"cycle":   c.Cycle,
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
	cycle := controls.NewObservableRatioDial()

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
		Cycle:   cycle,
		Speed:   speed,
		Chase:   &chase.FloatDial,
		Reverse: &reverse.Toggle,
	}
}
