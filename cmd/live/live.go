package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akualab/dmx"
	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
	"github.com/mholzen/play-go/stages"

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

	duration := clock.PhrasePeriod()
	transition := func() {
		start, end := seq.IncValues()

		// action := Transition(home.Universe, start, end, duration, ease.InOutSine, REFRESH)
		// action()

		// Chain
		log.Printf("loop with universe %+v", home.Universe)
		for i, f := range home.Universe {
			action := Transition(f, start, end, duration, ease.InOutSine, REFRESH)

			// d := time.Duration(i*100) * time.Millisecond // time.Duration(len(home.Universe))
			d := (duration / 2) * time.Duration(i)
			log.Printf("i: %d delay: %s", i, d)
			go Delay(d, action)
		}
	}

	clock.On(controls.TriggerOnBar(1), transition)
}

func moveTomshine() {
	current, _ := controls.NewMap("tilt:64-192", "pan:0-64")
	moveStep := func() {
		end, _ := controls.NewMap("tilt:64-192", "pan:0-64")
		Transition(home.TomeShine, current, end, clock.BeatPeriod(), ease.OutExpo, clock.BeatPeriod())()
		current = end
	}
	RepeatEvery(clock.PhrasePeriod(), moveStep)
}

func beatDown() {
	freedomPars := controls.NewSequenceT(home.FreedomPars)
	tomShines := controls.NewSequenceT(home.TomeShine)

	RepeatEvery(clock.BeatPeriod(), func() {
		freedomPar, _ := freedomPars.IncValues()
		tomShine, _ := tomShines.IncValues()

		duration := clock.BeatPeriod() / 2
		Transition(freedomPar, controls.ValueMap{"dimmer": 255}, controls.ValueMap{"dimmer": 0}, duration, ease.OutCubic, 10*time.Millisecond)()
		Transition(tomShine, controls.ValueMap{"dimmer": 255}, controls.ValueMap{"dimmer": 0}, duration, ease.OutCubic, 10*time.Millisecond)()

		freedomPar.SetValue("dimmer", 255)
		tomShine.SetValue("dimmer", 255)
	})
}

func setup() {
	controls.LoadColors()
	universe.SetAll(0)
	universe.SetValue("dimmer", 255)

	home.TomeShine.SetValue("tilt", 127)
	home.TomeShine.SetValue("speed", 255)

	home.ColorStrip.SetValue("mode", 210)

	clock.On(controls.TriggerOnBeats(), func() { log.Printf("clock: %s", clock.String()) })

	clock.Start()
}

func twoColors() {
	red := controls.AllColors["cyan"].Values()
	red.ApplyTo(home.FreedomPars.Odd())

	blue := controls.AllColors["yellow_green"].Values()
	blue.ApplyTo(home.FreedomPars.Even())
}

var clock = controls.NewClock(120)

func main() {
	connection, err := fixture.GetConnection()
	if err != nil {
		log.Printf("Warning: starting without a DMX connection: %s", err)
	}

	setup()

	rainbow()
	// twoColors()
	// gold()
	// moveTomshine()
	// beatDown()

	if connection != nil {
		fixture.Render(home.Universe, connection)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
}

// func GetConnection() *fixture.DMX {
// 	connection, err := fixture.NewDMXConnection("/dev/ttyUSB0") // Linux
// 	if err != nil {
// 		connection, err = fixture.NewDMXConnection("/dev/tty.usbserial-ENVVVCOF") // MacOS
// 		if err != nil {
// 			log.Printf("Warning: starting without a DMX connection: %s", err)
// 			connection = nil
// 		}
// 	}
// 	return connection
// }

const REFRESH = 40 * time.Millisecond // DMXIS cannot read faster than 40ms

func Render(f fixture.Fixtures, connection *dmx.DMX) {
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
