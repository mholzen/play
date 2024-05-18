package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
	"github.com/mholzen/play-go/stages"

	"github.com/akualab/dmx"
	"github.com/fogleman/ease"
)

var home = stages.GetHome()
var universe = home.Universe
var stepDuration = 5 * time.Second

func gold() {
	controls.AllColors["gold"].Values().ApplyTo(universe)
}

func rainbow() {
	seq := controls.NewSequence([]controls.ValueMap{
		controls.AllColors["red"].Values(), // TODO: if `red` doesn't exist, this should fail fast rather than return a the 0 (ie. black) color
		controls.AllColors["yellow"].Values(),
		controls.AllColors["green"].Values(),
		controls.AllColors["cyan"].Values(),
		controls.AllColors["blue"].Values(),
		controls.AllColors["purple"].Values(),
		controls.AllColors["violet"].Values(),
	})

	transition := func() {
		start, end := seq.IncValues()

		for i, f := range home.FreedomPars {
			action := Transition(f, start, end, stepDuration, ease.InOutSine, REFRESH)
			d := time.Duration(i) * stepDuration / time.Duration(len(home.FreedomPars))
			Delay(d, action)
		}
		// Transition(universe, start, end, stepDuration, ease.InOutSine, REFRESH*2)()
	}
	RepeatEvery(stepDuration, transition)
}

func moveTomshine() {
	current, _ := controls.NewMap("tilt:96-160", "pan:0-256")
	moveStep := func() {
		end, _ := controls.NewMap("tilt:96-160", "pan:0-256")
		Transition(home.TomeShine, current, end, stepDuration/4, ease.InOutSine, 10*time.Millisecond)()
		current = end
	}
	RepeatEvery(20*time.Second, moveStep)
}

func beatDown() {
	freedomPars := controls.NewSequenceT(home.FreedomPars)
	tomShines := controls.NewSequenceT(home.TomeShine)

	RepeatEvery(stepDuration/20, func() {
		freedomPar, _ := freedomPars.IncValues()
		tomShine, _ := tomShines.IncValues()

		Transition(freedomPar, controls.ValueMap{"dimmer": 255}, controls.ValueMap{"dimmer": 0}, stepDuration/20, ease.OutCubic, 10*time.Millisecond)()
		Transition(tomShine, controls.ValueMap{"dimmer": 255}, controls.ValueMap{"dimmer": 0}, stepDuration/20, ease.OutCubic, 10*time.Millisecond)()

		freedomPar.SetValue("dimmer", 255)
		tomShine.SetValue("dimmer", 255)
	})
}

func setup() {
	controls.LoadColors()
	universe.SetAll(0)
	universe.SetValue("dimmer", 24)
	home.ColorStrip.SetValue("mode", 210)

	home.TomeShine.SetValue("tilt", 127)
	home.TomeShine.SetValue("speed", 64)
}

func twoColors() {
	red := controls.AllColors["cyan"].Values()
	red.ApplyTo(home.FreedomPars.Odd())

	blue := controls.AllColors["yellow_green"].Values()
	blue.ApplyTo(home.FreedomPars.Even())
}

var clock = controls.NewClock(120)

func main() {
	connection, err := dmx.NewDMXConnection("/dev/ttyUSB0")
	if err != nil {
		log.Printf("Warning: starting without a DMX connection: %s", err)
		connection = nil
	}
	setup()

	rainbow()
	// twoColors()
	// gold()
	moveTomshine()
	// beatDown()

	if connection != nil {
		Render(home.Universe, connection)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
}

const REFRESH = 11 * time.Millisecond

func Render(f fixture.Fixtures2, connection *dmx.DMX) {
	ticker := time.NewTicker(REFRESH)
	go func() {
		for range ticker.C {
			f.Render(*connection)
		}
	}()
}

func Transition(f fixture.FixtureI, start, end controls.ValueMap, duration time.Duration, ease ease.Function, period time.Duration) func() {
	return func() {
		ticker := time.NewTicker(period)
		x := 0.0
		inc := float64(period.Milliseconds()) / float64(duration.Milliseconds())
		log.Printf("start: %s", start.String())
		// log.Printf("ending: %s", end.String())
		for range ticker.C {
			x += inc
			y := ease(x)

			// log.Printf("x: %f, y: %f", x, y)
			values := controls.InterpolateValues(start, end, y)
			// log.Printf("current: %s", values)
			values.ApplyTo(f)

			if x >= 1.0 {
				ticker.Stop()
				break
			}
		}
		log.Printf("end: %s", end.String())
	}
}

func RepeatEvery(duration time.Duration, action func()) {
	ticker := time.NewTicker(duration)
	go func() {
		action()
		for range ticker.C {
			action()
		}
	}()
}

func Delay(duration time.Duration, action func()) {
	time.Sleep(duration)
	action()
}
