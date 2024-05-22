//go:build linux

package controls

import (
	"fmt"
	"log"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/rtmididrv"
)

type TickerChannel <-chan int

var CLOCK = byte(0xF8)
var SYNC = byte(0xFA)

func GetMidiClockTicker() (TickerChannel, error) {

	// Initialize the MIDI driver
	drv, err := rtmididrv.New()
	if err != nil {
		log.Fatalf("could not start MIDI driver: %v", err)
	}
	defer drv.Close()

	// Get the list of available MIDI inputs
	inPorts, err := drv.Ins()
	if err != nil {
		log.Fatalf("could not get MIDI input ports: %v", err)
	}

	if len(inPorts) == 0 {
		log.Fatalf("no MIDI input ports available")
	}

	// Open the first available MIDI input port
	in := inPorts[0]
	err = in.Open()
	if err != nil {
		log.Fatalf("could not open MIDI input port: %v", err)
	}
	defer in.Close()

	fmt.Printf("Listening on MIDI input: %s\n", in.String())

	// Create a MIDI reader
	tick := 0
	ch := make(chan int)

	r := reader.New(
		reader.NoLogger(),
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			raw := msg.Raw()
			if len(raw) > 0 && raw[0] == 0xF8 {
				tick++
				ch <- tick
			}
		}),
	)

	// Listen for MIDI messages
	go r.ListenTo(in)

	// Keep the program running
	select {}
}
