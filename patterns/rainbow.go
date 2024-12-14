package patterns

import (
	"time"

	"github.com/fogleman/ease"
	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
)

func Rainbow(fixtures *fixture.FixturesGeneric[fixture.FixtureI], clock *controls.Clock) controls.Triggers {
	seq := controls.NewSequence([]controls.ChannelValues{
		controls.AllColors["red"].Values(), // TODO: if `red` doesn't exist, this should fail fast rather than return a the 0 (ie. black) color
		controls.AllColors["yellow"].Values(),
		controls.AllColors["green"].Values(),
		controls.AllColors["cyan"].Values(),
		controls.AllColors["blue"].Values(),
		controls.AllColors["purple"].Values(),
	})

	duration := clock.PhrasePeriod()
	transition := func() {
		start, end := seq.IncValues()
		// log.Printf("transition %v to %v\n", start, end)
		for i, f := range fixtures.GetFixtureList() {
			action := Transition(f, start, end, duration, ease.InOutSine, fixture.REFRESH)

			f.SetChannelValue("dimmer", 255)
			f.SetChannelValue("tilt", 127)

			d := time.Duration(i) * (duration / 2)
			go Delay(d, action)
		}
	}

	t := clock.On(controls.TriggerOnBar(1), transition)
	return controls.Triggers{
		*t,
	}
}
