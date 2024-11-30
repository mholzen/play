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
	err := controls.LoadColors()
	if err != nil {
		log.Fatalf("Error loading colors: %v", err)
	}

	var home = stages.GetHome()
	var universe = home.Universe

	universe.SetAll(0)
	universe.SetChannelValue("dimmer", 0)
	universe.SetChannelValue("mode", 210) // for colorstrip mini

	soft_white := controls.AllColors["soft_white"]
	universe.SetValueMap(soft_white.Values())

	clock := controls.NewClock(120)
	clock.On(controls.TriggerOnBeats(), func() { log.Printf("clock: %s", clock.String()) })
	clock.Start()

	surface := GetRootSurface(universe, clock)

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

func NewDialMapAllFixtures(fixtures fixture.Fixtures) controls.DialMap {
	controls := NewDialMap()
	go func() {
		channel := controls.Channel()
		for {
			valueMap := <-channel
			fixtures.SetValueMap(valueMap)
		}
	}()
	return controls
}

func LinkEmitterToFixture(source controls.Emitter[fixture.FixtureValues], target fixture.Fixtures) {
	go func() {
		for fixtureValues := range source.Channel() {
			log.Printf("setting fixture values: %v", fixtureValues)
			target.SetValue(fixtureValues)
		}
	}()
}

func GetRootSurface(universe fixture.Fixtures, clock *controls.Clock) controls.Container {
	surface := controls.NewList(3)

	dialFixtures := fixture.NewFixturesFromFixtures(universe)
	channelDials := NewDialMapAllFixtures(dialFixtures)
	surface.SetItem(0, channelDials)

	rainbowFixtures := fixture.NewFixturesFromFixtures(universe)
	patterns.Rainbow(rainbowFixtures, clock)

	mux := controls.NewMux[fixture.FixtureValues]()
	mux.Add("dials", dialFixtures)
	mux.Add("rainbow", rainbowFixtures)
	surface.SetItem(2, &mux)

	// link mux emitter to universe fixture
	LinkEmitterToFixture(mux, universe)

	return surface
}

func NewDialMap() controls.DialMap {
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
