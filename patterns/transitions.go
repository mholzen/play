package patterns

import (
	"log"
	"time"

	"github.com/mholzen/play-go/controls"
)

func Repeat(duration time.Duration, f func()) *time.Ticker {
	f()
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			f()
		}
	}()
	return ticker
}

func Ease(dial *controls.NumericDial, duration time.Duration, endValue byte) {
	startValue := float64(dial.Value)
	startTime := time.Now()
	ticker := time.NewTicker(REFRESH).C
	go func() {
		for t := range ticker {
			if t.After(startTime.Add(duration)) {
				dial.SetValue(endValue)
				return
			}
			elapsed := t.Sub(startTime)
			step := float64(elapsed) / float64(duration)
			// mult := ease.InOutSine(step)
			mult := step
			value := startValue + (float64(endValue)-startValue)*mult

			log.Printf("elapsed: %f step: %f mult: %f value: %f",
				float64(elapsed), step, mult, value,
			)

			dial.SetValue(byte(value))
		}
	}()
}
