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

	// surface := GetChannelMap(universe)
	surface := GetRootSurface(universe)

	// Repeat(8*time.Second, GetToggleFunc(dimmerDial, 6*time.Second))

	if connection != nil {
		fixture.Render(universe, connection)
	}

	StartServer(surface)

	time.Sleep(100 * time.Second)
}

func GetChannelMap(universe fixture.Fixtures) controls.DialMap {
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
	return surface
}

func GetRootSurface(universe fixture.Fixtures) controls.Container {
	surface := controls.NewList(2)
	surface.SetItem(0, GetChannelMap(universe))
	surface.SetItem(1, GetToggles(universe))
	return surface
}

func NewControls() controls.DialMap {
	m := make(controls.DialMap)
	m["mode"] = controls.NewNumericDial()

	m["dimmer"] = controls.NewNumericDial()
	m["strobe"] = controls.NewNumericDial()

	m["tilt"] = controls.NewNumericDial()
	m["pan"] = controls.NewNumericDial()
	m["speed"] = controls.NewNumericDial()

	m["r"] = controls.NewNumericDial()
	m["g"] = controls.NewNumericDial()
	m["b"] = controls.NewNumericDial()
	m["w"] = controls.NewNumericDial()
	m["a"] = controls.NewNumericDial()
	m["uv"] = controls.NewNumericDial()

	return m
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
