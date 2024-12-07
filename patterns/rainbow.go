package patterns

import (
	"log"
	"time"

	"github.com/fogleman/ease"
	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
)

func Rainbow(fixtures *fixture.FixturesGeneric[fixture.FixtureI], clock *controls.Clock) controls.Triggers {
	seq := controls.NewSequence([]controls.ValueMap{
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
		log.Printf("transition %v to %v\n", start, end)
		for i, f := range fixtures.GetFixtures() {
			action := Transition(f, start, end, duration, ease.InOutSine, fixture.REFRESH)

			d := (duration / 2) * time.Duration(i)
			go Delay(d, action)
		}
	}

	t := clock.On(controls.TriggerOnBar(1), transition)
	return controls.Triggers{
		*t,
	}
}
