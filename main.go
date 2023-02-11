package main

import (
	"log"
	"play-go/controls"
	"play-go/fixture"
	"time"

	"github.com/akualab/dmx"
)

func main() {
	connection, err := dmx.NewDMXConnection("/dev/ttyUSB0")
	if err != nil {
		log.Fatal(err)
	}

	universe := fixture.NewFixtureList()
	addresses := []int{65, 81, 97, 113}
	for _, address := range addresses {
		universe.AddFixture(fixture.NewFreedomPar(), address)
	}

	soft_white := controls.AllColors["soft_white"]

	universe.SetAll(0)
	universe.SetValue("dimmer", 0)

	soft_white.Values().ApplyTo(universe)

	universe.Render(*connection)

	surface := GetControls()
	dimmerDial := surface["dimmer"]
	controls.LinkFixtureChannel(universe, "dimmer", dimmerDial.Channel)
	Render(universe, connection)
	Repeat(8*time.Second, GetToggleFunc(dimmerDial, 6*time.Second))

	time.Sleep(100 * time.Second)
}

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

func GetToggleFunc(dial *controls.Dial, duration time.Duration) func() {
	return func() {
		endValue := dial.Opposite()
		log.Printf("triggering ease from %d to %d", dial.Value, endValue)
		Ease(dial, duration/2, endValue)
	}
}

func Ease(dial *controls.Dial, duration time.Duration, endValue byte) {
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

func GetControls() controls.DialMap {
	m := make(controls.DialMap)
	m["dimmer"] = controls.NewDial()
	return m
}

const REFRESH = 11 * time.Millisecond

func Render(f fixture.FixtureList, connection *dmx.DMX) {
	ticker := time.NewTicker(REFRESH)
	go func() {
		for range ticker.C {
			f.Render(*connection)
		}
	}()
}
