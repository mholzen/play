package home

import (
	"context"
	"time"

	"github.com/fogleman/ease"
	"github.com/mholzen/play/controls"
	"github.com/mholzen/play/fixture"
	"github.com/mholzen/play/patterns"
)

var home = GetHome()

var Transition = patterns.Transition
var Delay = patterns.Delay
var RepeatEvery = patterns.RepeatEvery

var clock = controls.NewClock(120)

func moveTomshine() controls.Triggers {
	current, _ := controls.NewChannelValues("tilt:64-192", "pan:0-64")
	moveStep := func() {
		end, _ := controls.NewChannelValues("tilt:64-192", "pan:0-64")
		Transition(home.TomeShine, current, end, clock.BeatPeriod(), ease.OutExpo, clock.BeatPeriod(), context.Background())()
		current = end
	}
	return controls.Triggers{
		*RepeatEvery(clock.PhrasePeriod(), moveStep),
	}
}

func beatDown() controls.Triggers {
	freedomPars := controls.NewSequence(home.FreedomPars.GetFixtureList())
	tomShines := controls.NewSequence(home.TomeShine.GetFixtureList())

	return controls.Triggers{
		*RepeatEvery(clock.BeatPeriod(), func() {
			freedomPar, _ := freedomPars.IncValues()
			tomShine, _ := tomShines.IncValues()

			duration := clock.BeatPeriod()
			Transition(home.FreedomPars, controls.ChannelValues{"dimmer": 255}, controls.ChannelValues{"dimmer": 0}, duration, ease.OutCubic, 10*time.Millisecond, context.Background())()
			Transition(tomShine, controls.ChannelValues{"dimmer": 255}, controls.ChannelValues{"dimmer": 0}, duration, ease.OutCubic, 10*time.Millisecond, context.Background())()

			freedomPar.SetChannelValue("dimmer", 255)
			tomShine.SetChannelValue("dimmer", 255)
		}),
	}
}

func moveDownTomshine() controls.Triggers {
	top, _ := controls.NewChannelValues("tilt:128", "pan:255")
	bottom, _ := controls.NewChannelValues("tilt:0", "pan:255")
	tomShines := controls.NewSequence(home.TomeShine.GetFixtureList())

	home.TomeShine.SetChannelValue("speed", 0)
	return controls.Triggers{
		*RepeatEvery(clock.BarPeriod(), func() {

			tomShine, _ := tomShines.IncValues()

			tomShine.SetChannelValue("tilt", 128)
			// tomShine.SetValue("dimmer", 255)
			Transition(tomShine, top, bottom, clock.BeatPeriod(), ease.Linear, 10*time.Millisecond, context.Background())()
			time.Sleep(clock.BeatPeriod())
			// tomShine.SetValue("dimmer", 0)
			tomShine.SetChannelValue("tilt", 128)
		}),
	}
}

func twoColors() controls.Triggers {
	colorA := controls.AllColors["cyan"].Values()
	colorB := controls.AllColors["yellow_green"].Values()

	return controls.Triggers{
		*RepeatEvery(clock.BeatPeriod(), func() {
			colorA, colorB = colorB, colorA
			fixture.ApplyTo(colorA, home.FreedomPars.Odd())
			fixture.ApplyTo(colorB, home.FreedomPars.Even())
		}),
	}
}

func GetTransitions() map[string]controls.Triggers {
	return map[string]controls.Triggers{
		// "rainbow":          Rainbow(),
		"beatDown":         beatDown(),
		"moveTomshine":     moveTomshine(),
		"moveDownTomshine": moveDownTomshine(),
		"twoColors":        twoColors(),
	}
}
