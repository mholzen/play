package patterns

import (
	"time"

	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/stages"

	"github.com/fogleman/ease"
)

var Transition = controls.Transition
var Delay = controls.Delay
var RepeatEvery = controls.RepeatEvery

var home = stages.GetHome()
var universe = home.Universe

func rainbow() controls.Triggers {
	seq := controls.NewSequence([]controls.ValueMap{
		controls.AllColors["red"].Values(), // TODO: if `red` doesn't exist, this should fail fast rather than return a the 0 (ie. black) color
		controls.AllColors["yellow"].Values(),
		controls.AllColors["green"].Values(),
		controls.AllColors["cyan"].Values(),
		controls.AllColors["blue"].Values(),
		controls.AllColors["purple"].Values(),
	})

	duration := clock.BarPeriod()
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
	return controls.Triggers{
		*t,
	}
}

func moveTomshine() controls.Triggers {
	current, _ := controls.NewMap("tilt:64-192", "pan:0-64")
	moveStep := func() {
		end, _ := controls.NewMap("tilt:64-192", "pan:0-64")
		Transition(home.TomeShine, current, end, clock.BeatPeriod(), ease.OutExpo, clock.BeatPeriod())()
		current = end
	}
	return controls.Triggers{
		*RepeatEvery(clock.PhrasePeriod(), moveStep),
	}
}

func beatDown() controls.Triggers {
	freedomPars := controls.NewSequenceT(home.FreedomPars)
	tomShines := controls.NewSequenceT(home.TomeShine)

	return controls.Triggers{
		*RepeatEvery(clock.BeatPeriod(), func() {
			freedomPar, _ := freedomPars.IncValues()
			tomShine, _ := tomShines.IncValues()

			duration := clock.BeatPeriod()
			Transition(home.FreedomPars, controls.ValueMap{"dimmer": 255}, controls.ValueMap{"dimmer": 0}, duration, ease.OutCubic, 10*time.Millisecond)()
			Transition(tomShine, controls.ValueMap{"dimmer": 255}, controls.ValueMap{"dimmer": 0}, duration, ease.OutCubic, 10*time.Millisecond)()

			freedomPar.SetValue("dimmer", 255)
			tomShine.SetValue("dimmer", 255)
		}),
	}
}

func moveDownTomshine() controls.Triggers {
	top, _ := controls.NewMap("tilt:128", "pan:255")
	bottom, _ := controls.NewMap("tilt:0", "pan:255")
	tomShines := controls.NewSequenceT(home.TomeShine)

	home.TomeShine.SetValue("speed", 0)
	return controls.Triggers{
		*RepeatEvery(clock.BarPeriod(), func() {

			tomShine, _ := tomShines.IncValues()

			tomShine.SetValue("tilt", 128)
			// tomShine.SetValue("dimmer", 255)
			Transition(tomShine, top, bottom, clock.BeatPeriod(), ease.Linear, 10*time.Millisecond)()
			time.Sleep(clock.BeatPeriod())
			// tomShine.SetValue("dimmer", 0)
			tomShine.SetValue("tilt", 128)
		}),
	}
}

func twoColors() {
	red := controls.AllColors["cyan"].Values()
	red.ApplyTo(home.FreedomPars.Odd())

	blue := controls.AllColors["yellow_green"].Values()
	blue.ApplyTo(home.FreedomPars.Even())
}

var clock = controls.NewClock(120)

const REFRESH = 40 * time.Millisecond // DMXIS cannot read faster than 40ms

func GetTransitions() map[string]controls.Triggers {
	return map[string]controls.Triggers{
		"rainbow":          rainbow(),
		"beatDown":         beatDown(),
		"moveTomshine":     moveTomshine(),
		"moveDownTomshine": moveDownTomshine(),
	}
}
