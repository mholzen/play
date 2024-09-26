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
	"github.com/nsf/termbox-go"

	"github.com/fogleman/ease"
)

var Transition = controls.Transition
var Delay = controls.Delay
var RepeatEvery = controls.RepeatEvery

var home = stages.GetHome()
var universe = home.Universe

func setup() {
	controls.LoadColors()
	universe.SetAll(0)
	universe.SetValue("dimmer", 32)

	home.TomeShine.SetValue("tilt", 127)
	home.TomeShine.SetValue("speed", 255)

	home.ColorStrip.SetValue("mode", 210)

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

func gold() {
	controls.AllColors["gold"].Values().ApplyTo(universe)
}

func rainbow() []*controls.Trigger {
	seq := controls.NewSequence([]controls.ValueMap{
		controls.AllColors["red"].Values(), // TODO: if `red` doesn't exist, this should fail fast rather than return a the 0 (ie. black) color
		controls.AllColors["yellow"].Values(),
		controls.AllColors["green"].Values(),
		controls.AllColors["cyan"].Values(),
		controls.AllColors["blue"].Values(),
		controls.AllColors["purple"].Values(),
		// controls.AllColors["violet"].Values(),
	})

	duration := clock.PhrasePeriod()
	transition := func() {
		start, end := seq.IncValues()

		// action := Transition(home.Universe, start, end, duration, ease.InOutSine, REFRESH)
		// action()

		// Chain
		// log.Printf("loop with universe %+v", home.Universe)
		for i, f := range home.Universe {
			action := Transition(f, start, end, duration, ease.InOutSine, REFRESH)

			// d := time.Duration(i*100) * time.Millisecond // time.Duration(len(home.Universe))
			d := (duration / 2) * time.Duration(i)
			// log.Printf("i: %d delay: %s", i, d)
			go Delay(d, action)
		}
	}

	t := clock.On(controls.TriggerOnBar(1), transition)
	return []*controls.Trigger{t}
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

		duration := clock.BeatPeriod()
		Transition(home.FreedomPars, controls.ValueMap{"dimmer": 255}, controls.ValueMap{"dimmer": 0}, duration, ease.OutCubic, 10*time.Millisecond)()
		Transition(tomShine, controls.ValueMap{"dimmer": 255}, controls.ValueMap{"dimmer": 0}, duration, ease.OutCubic, 10*time.Millisecond)()

		freedomPar.SetValue("dimmer", 255)
		tomShine.SetValue("dimmer", 255)
	})
}

func moveDownTomshine() {
	top, _ := controls.NewMap("tilt:128", "pan:255")
	bottom, _ := controls.NewMap("tilt:0", "pan:255")
	tomShines := controls.NewSequenceT(home.TomeShine)

	home.TomeShine.SetValue("speed", 0)
	RepeatEvery(clock.BarPeriod(), func() {

		tomShine, _ := tomShines.IncValues()

		tomShine.SetValue("tilt", 128)
		// tomShine.SetValue("dimmer", 255)
		Transition(tomShine, top, bottom, clock.BeatPeriod(), ease.Linear, 10*time.Millisecond)()
		time.Sleep(clock.BeatPeriod())
		// tomShine.SetValue("dimmer", 0)
		tomShine.SetValue("tilt", 128)
	})
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
	// moveDownTomshine()
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

const REFRESH = 40 * time.Millisecond // DMXIS cannot read faster than 40ms

func GetSymbols() []string {
	return []string{
		"Transition",
		"Delay",
		"RepeatEvery",
		"home",
		"universe",
		"setup",
		"gold",
		"rainbow",
		"moveTomshine",
		"beatDown",
		"moveDownTomshine",
		"twoColors",
		"clock",
		"main",
		"REFRESH",
	}
}
