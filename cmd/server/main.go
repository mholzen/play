package main

import (
	"log"

	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
	"github.com/mholzen/play-go/stages/home"
)

func main() {

	var h = home.GetHome()
	var universe = h.Universe

	universe.SetAll(0)
	universe.SetChannelValue("dimmer", 0)
	universe.SetChannelValue("mode", 210) // for colorstrip mini

	soft_white := controls.AllColors["soft_white"]
	universe.SetChannelValues(soft_white.Values())

	clock := controls.NewClock(120)
	clock.On(controls.TriggerOnBeats(), func() { log.Printf("clock: %s", clock.String()) })
	clock.Start()

	surface := home.GetRootSurface(universe, clock)

	connection, err := fixture.GetConnection()
	if err != nil {
		log.Printf("Warning: starting without a DMX connection: %s", err)
	}
	if connection != nil {
		fixture.Render(universe, *connection)
	}

	StartServer(surface)
}
