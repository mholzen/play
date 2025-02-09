package patterns

import (
	"time"

	"github.com/fogleman/ease"
	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
)

func Transition(f fixture.Fixture, start, end controls.ChannelValues, duration time.Duration, ease ease.Function, period time.Duration) func() {
	return TransitionValues(start, end, duration, ease, period, func(values controls.ChannelValues) {
		fixture.ApplyTo(values, f)
	})
}

func TransitionValues(start, end controls.ChannelValues, duration time.Duration, ease ease.Function, period time.Duration, apply func(values controls.ChannelValues)) func() {
	return func() {
		ticker := time.NewTicker(period)
		x := 0.0
		inc := float64(period) / float64(duration)

		// log.Printf("start: %s", start.String())
		// log.Printf("ending: %s", end.String())
		for range ticker.C {
			x += inc
			y := ease(x)

			// log.Printf("x: %f, y: %f", x, y)
			values := controls.InterpolateValues(start, end, y)
			// log.Printf("current: %s", values)
			apply(values)
			if x >= 1.0 {
				ticker.Stop()
				break
			}
		}
		// log.Printf("end: %s", end.String())
	}
}

func RepeatEvery(duration time.Duration, action func()) *controls.Trigger {
	trigger := &controls.Trigger{
		Enabled: true,
		Do:      action,
	}

	ticker := time.NewTicker(duration)
	go func() {
		action()
		for range ticker.C {
			if trigger.Enabled {
				action()
			}
		}
	}()
	return trigger
}

func Delay(duration time.Duration, action func()) {
	time.Sleep(duration)
	action()
}
