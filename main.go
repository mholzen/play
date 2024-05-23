package main

import (
	"log"
	"time"

	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
	"github.com/mholzen/play-go/stages"
)

func main() {
	connection, err := fixture.GetConnection()
	if err != nil {
		log.Printf("Warning: starting without a DMX connection: %s", err)
	}

	var home = stages.GetHome()
	var universe = home.Universe
	// universe := fixture.GetUniverse()

	universe.SetAll(0)
	universe.SetValue("dimmer", 0)
	universe.SetValue("mode", 210) // for colorstrip mini

	soft_white := controls.AllColors["soft_white"]
	soft_white.Values().ApplyTo(universe)

	surface := NewControls()

	controls.LinkDialToFixtureChannel(surface["mode"], universe, "mode")

	controls.LinkDialToFixtureChannel(surface["dimmer"], universe, "dimmer")
	controls.LinkDialToFixtureChannel(surface["strobe"], universe, "strobe")

	controls.LinkDialToFixtureChannel(surface["tilt"], universe, "tilt")
	controls.LinkDialToFixtureChannel(surface["pan"], universe, "pan")
	controls.LinkDialToFixtureChannel(surface["speed"], universe, "speed")

	controls.LinkDialToFixtureChannel(surface["r"], universe, "r")
	controls.LinkDialToFixtureChannel(surface["g"], universe, "g")
	controls.LinkDialToFixtureChannel(surface["b"], universe, "b")
	controls.LinkDialToFixtureChannel(surface["w"], universe, "w")
	controls.LinkDialToFixtureChannel(surface["a"], universe, "a")
	controls.LinkDialToFixtureChannel(surface["uv"], universe, "uv")

	// Repeat(8*time.Second, GetToggleFunc(dimmerDial, 6*time.Second))

	if connection != nil {
		fixture.Render(universe, connection)
	}

	StartServer(surface)

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

func NewControls() controls.DialMap {
	m := make(controls.DialMap)
	m["mode"] = controls.NewDial()

	m["dimmer"] = controls.NewDial()
	m["strobe"] = controls.NewDial()

	m["tilt"] = controls.NewDial()
	m["pan"] = controls.NewDial()
	m["speed"] = controls.NewDial()

	m["r"] = controls.NewDial()
	m["g"] = controls.NewDial()
	m["b"] = controls.NewDial()
	m["w"] = controls.NewDial()
	m["a"] = controls.NewDial()
	m["uv"] = controls.NewDial()
	return m
}

const REFRESH = 40 * time.Millisecond // DMXIS cannot read faster than 40ms
