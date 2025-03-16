package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
	"github.com/mholzen/play-go/patterns"
	"github.com/mholzen/play-go/stages/home"
)

var Transition = patterns.Transition
var Delay = patterns.Delay
var RepeatEvery = patterns.RepeatEvery

var h = home.GetHome()
var universe = h.Universe
var clock = controls.NewClock(120)
var termTrigger *controls.TermTrigger
var loggingEnabled = true // Track whether beat logging is enabled

func setup() {
	controls.LoadColors()
	universe.SetAll(0)
	universe.SetChannelValue("dimmer", 32)

	h.TomeShine.SetChannelValue("tilt", 127)
	h.TomeShine.SetChannelValue("speed", 255)

	h.ColorStrip.SetChannelValue("mode", 210)

	frequency := controls.NewTimeKeeper(10)

	termTrigger = controls.NewTerminalWithClock(controls.ClockControls{
		Clock:           clock,
		FrequencyKeeper: frequency,
		LoggingEnabled:  &loggingEnabled,
	})

	termTrigger.Start()

	clock.On(controls.TriggerOnBeats(), func() {
		if loggingEnabled {
			log.Printf("clock: %s", clock.String())
		}
	})
	clock.Start()
}

func main() {
	connection, err := fixture.GetConnection()
	if err != nil {
		log.Printf("Warning: starting without a DMX connection: %s", err)
	}

	setup()

	home.Rainbow()
	// moveDownTomshine()
	// twoColors()
	// gold()
	// moveTomshine()
	// beatDown()

	if connection != nil {
		fixture.Render(h.Universe, *connection)
	}

	// Set up signal handling for proper cleanup
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	<-sigs

	// Ensure terminal is reset before exiting
	if termTrigger != nil {
		termTrigger.Stop()
	}

	log.Println("Program exited cleanly")
}

const REFRESH = 40 * time.Millisecond // DMXIS cannot read faster than 40ms
