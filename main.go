package main

import (
	"log"
	"time"

	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
	"github.com/mholzen/play-go/patterns"
	"github.com/mholzen/play-go/stages"
)

func main() {

	var home = stages.GetHome()
	var universe = home.Universe

	universe.SetAll(0)
	universe.SetChannelValue("dimmer", 0)
	universe.SetChannelValue("mode", 210) // for colorstrip mini

	soft_white := controls.AllColors["soft_white"]
	fixture.ApplyTo(soft_white.Values(), universe)

	clock := controls.NewClock(120)

	surface := GetRootSurface(universe, *clock)

	connection, err := fixture.GetConnection()
	if err != nil {
		log.Printf("Warning: starting without a DMX connection: %s", err)
	}
	if connection != nil {
		fixture.Render(universe, connection)
	}

	StartServer(surface)

	time.Sleep(100 * time.Second)
}

func GetChannelMap(universe fixture.Fixtures) controls.DialMap {
	surface := NewControls()

	fixture.LinkDialToFixtureChannel(surface.Dials["mode"], universe, "mode")

	fixture.LinkDialToFixtureChannel(surface.Dials["dimmer"], universe, "dimmer")
	fixture.LinkDialToFixtureChannel(surface.Dials["strobe"], universe, "strobe")

	fixture.LinkDialToFixtureChannel(surface.Dials["tilt"], universe, "tilt")
	fixture.LinkDialToFixtureChannel(surface.Dials["pan"], universe, "pan")
	fixture.LinkDialToFixtureChannel(surface.Dials["speed"], universe, "speed")

	fixture.LinkDialToFixtureChannel(surface.Dials["r"], universe, "r")
	fixture.LinkDialToFixtureChannel(surface.Dials["g"], universe, "g")
	fixture.LinkDialToFixtureChannel(surface.Dials["b"], universe, "b")
	fixture.LinkDialToFixtureChannel(surface.Dials["w"], universe, "w")
	fixture.LinkDialToFixtureChannel(surface.Dials["a"], universe, "a")
	fixture.LinkDialToFixtureChannel(surface.Dials["uv"], universe, "uv")
	return surface
}

func GetRootSurface(universe fixture.Fixtures, clock controls.Clock) controls.Container {
	surface := controls.NewList(3)
	surface.SetItem(0, GetChannelMap(universe))

	rainbowFixtures := fixture.NewFixturesFromFixtures(universe)
	surface.SetItem(1, patterns.Rainbow(rainbowFixtures, clock))

	mux := controls.NewMux[fixture.FixtureValues]()
	mux.Add("rainbow", rainbowFixtures)
	mux.Add("dials", universe)
	surface.SetItem(2, mux)

	return surface
}

func NewControls() controls.DialMap {
	return controls.NewNumericDialMap("mode", "dimmer", "strobe", "tilt", "pan", "speed", "r", "g", "b", "w", "a", "uv")
}

func LinkToggleToEnable(toggle *controls.Toggle, recipient controls.Triggers) {
	go func() {
		for value := range toggle.C {
			if value {
				recipient.Enable()
			} else {
				recipient.Disable()
			}
		}
	}()
}

func GetToggles(universe fixture.Fixtures) *controls.Toggle {
	t1 := controls.NewToggle()
	transitions := patterns.GetTransitions()
	LinkToggleToEnable(t1, transitions["rainbow"])
	return t1
}

const REFRESH = 40 * time.Millisecond // DMXIS cannot read faster than 40ms
