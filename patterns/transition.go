package patterns

import (
	"context"
	"time"

	"github.com/fogleman/ease"
	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
)

func Transition(f fixture.Fixture, start, end controls.ChannelValues, duration time.Duration, ease ease.Function, period time.Duration, ctx context.Context) func() {
	return TransitionValues(start, end, duration, ease, period, func(values controls.ChannelValues) {
		fixture.ApplyTo(values, f)
	}, ctx)
}

func TransitionValues(start, end controls.ChannelValues, duration time.Duration, ease ease.Function, period time.Duration, apply func(values controls.ChannelValues), ctx context.Context) func() {
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

			select {
			case <-ctx.Done():
				x = 1.0
			default:
			}

			if x >= 1.0 {
				ticker.Stop()

				// Ensure we apply the final values.
				// This could be made optional in the future.
				// InterpolateValues above does not always end with the final value.
				apply(end)

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

func Delay(duration time.Duration, action func(), ctx context.Context) {
	select {
	case <-ctx.Done():
		// log.Printf("delay cancelled")
		return
	case <-time.After(duration):
		// log.Printf("executing action after delay: %v", duration)
		action()
	}
}
