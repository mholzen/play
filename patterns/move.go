package patterns

import (
	"context"
	"fmt"
	"time"

	"github.com/fogleman/ease"
	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
)

type MoveControls struct {
	// controls are more than just the configurable parameters.  it also contains the dials, knobs and buttons used to select the parameters.
	Clock      *controls.Clock               `json:"-"`
	Speed      *controls.ObservableRatioDial `json:"speed"` // TODO: doesn't have to be observable
	Randomness *controls.ObservableFloatDial `json:"randomness"`
}

func (c MoveControls) Move(fixtures *fixture.AddressableFixtures[fixture.Fixture]) controls.Triggers {
	transition := func() {
		ratio := c.Speed.Get().ToFloat()
		duration := time.Duration(float64(c.Clock.PhrasePeriod().Nanoseconds()) * ratio)

		max := len(fixtures.GetFixtureList())
		contexts := make([]context.Context, max)
		cancels := make([]context.CancelFunc, max)

		for i, f := range fixtures.GetFixtureList() {
			f.SetChannelValue("tilt", byte(randomTilt()))
			f.SetChannelValue("pan", byte(randomPan()))

			if cancels[i] != nil {
				cancels[i]()
			}
			contexts[i], cancels[i] = context.WithCancel(context.Background())
			action := Transition(f, controls.ChannelValues{"tilt": byte(randomTilt()), "pan": byte(randomPan())}, controls.ChannelValues{"tilt": byte(randomTilt()), "pan": byte(randomPan())}, duration, ease.InOutSine, fixture.REFRESH, contexts[i])

			go Delay(time.Duration(0), action, contexts[i])
		}
	}

	t := c.Clock.On(controls.TriggerOnBar(1), transition)
	return controls.Triggers{
		*t,
	}
}

func (c MoveControls) Items() map[string]controls.Item {
	return map[string]controls.Item{
		"speed":      c.Speed,
		"randomness": c.Randomness,
	}
}

func (c MoveControls) GetItem(name string) (controls.Item, error) {
	items := c.Items()
	if item, ok := items[name]; ok {
		return item, nil
	}
	return nil, fmt.Errorf("item not found")
}

func NewMoveControls(clock *controls.Clock) MoveControls {
	speed := controls.NewObservableRatioDial()
	randomness := NewObservableFloatDial()

	return MoveControls{
		Clock:      clock,
		Speed:      speed,
		Randomness: randomness,
	}
}

func randomTilt() int {
	// Implement logic to generate a random tilt value
	return 127 // Placeholder value
}

func randomPan() int {
	// Implement logic to generate a random pan value
	return 127 // Placeholder value
}

func NewObservableFloatDial() *controls.ObservableFloatDial {
	return &controls.ObservableFloatDial{
		FloatDial: controls.FloatDial{
			Value: 0.5,
			Min:   0.0,
			Max:   1.0,
		},
	}
}
