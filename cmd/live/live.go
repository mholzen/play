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
	"github.com/nsf/termbox-go"
)

var Transition = patterns.Transition
var Delay = patterns.Delay
var RepeatEvery = patterns.RepeatEvery

var h = home.GetHome()
var universe = h.Universe

func setup() {
	controls.LoadColors()
	universe.SetAll(0)
	universe.SetChannelValue("dimmer", 32)

	h.TomeShine.SetChannelValue("tilt", 127)
	h.TomeShine.SetChannelValue("speed", 255)

	h.ColorStrip.SetChannelValue("mode", 210)

	frequency := controls.NewTimeKeeper(10)

	controls.NewTermTrigger(
		map[termbox.Key]func(){
			termbox.KeyCtrlC: func() {
				log.Printf("quitting")
				os.Exit(0)
			},
			termbox.Key('t'): func() {
				frequency.AddTime(time.Now())
				bpm, _ := frequency.GetBpm()
				log.Printf("=== TAP bpm: %f", bpm)
			},
			termbox.Key('s'): func() {
				bpm, _ := frequency.GetBpm()
				clock.SetBpm(bpm)
				log.Printf("=== SYNC bpm: %f", bpm)
			},
			termbox.Key('r'): func() {
				clock.Reset()
				log.Printf("=== RESET")
			},
			termbox.Key('['): func() {
				clock.Nudge(-10 * time.Millisecond)
				log.Printf("=== NUDGE BACK -10ms")
			},
			termbox.Key(']'): func() {
				clock.Nudge(10 * time.Millisecond)
				log.Printf("=== NUDGE FORWARD +10ms")
			},
			termbox.Key('+'): func() {
				clock.Pace(.01)
				log.Printf("=== PACE UP .01")
			},
			termbox.Key('-'): func() {
				clock.Pace(-.01)
				log.Printf("=== PACE DOWN .01")
			},
		},
		map[termbox.Key]func(float64){
			termbox.Key('m'): func(bpm float64) {
				clock.SetBpm(bpm)
				log.Printf("=== SET bpm: %f", bpm)
			},
		},
	).Start()

	clock.On(controls.TriggerOnBeats(), func() { log.Printf("clock: %s", clock.String()) })
	clock.Start()
}

var clock = controls.NewClock(120)

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
		fixture.Render(h.Universe, connection)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
}

const REFRESH = 40 * time.Millisecond // DMXIS cannot read faster than 40ms
