package main

import (
	"log"
	"time"

	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
	"github.com/mholzen/play-go/stages/home"
)

func main() {
	err := controls.LoadColors()
	if err != nil {
		log.Fatalf("Error loading colors: %v", err)
	}

	var h = home.GetHome()
	var universe = h.Universe

	universe.SetAll(0)
	universe.SetChannelValue("dimmer", 0)
	universe.SetChannelValue("mode", 210) // for colorstrip mini

	soft_white := controls.AllColors["soft_white"]
	universe.SetValueMap(soft_white.Values())

	clock := controls.NewClock(120)
	clock.On(controls.TriggerOnBars(), func() { log.Printf("clock: %s", clock.String()) })
	clock.Start()

	surface := home.GetRootSurface(universe, clock)

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

const REFRESH = 40 * time.Millisecond // DMXIS cannot read faster than 40ms
