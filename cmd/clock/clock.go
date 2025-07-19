package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mholzen/play/controls"
)

var clock = controls.NewClock(120)
var termTrigger *controls.TermTrigger
var loggingEnabled = true

func setup() {
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
	setup()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	if termTrigger != nil {
		termTrigger.Stop()
	}

	log.Println("Program exited cleanly")
}
