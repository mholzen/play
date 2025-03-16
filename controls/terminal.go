package controls

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nsf/termbox-go"
)

type ClockControls struct {
	Clock           *Clock
	FrequencyKeeper *TimeKeeper
	LoggingEnabled  *bool
	OnQuit          func()
}

func NewTerminalWithClock(controls ClockControls) *TermTrigger {
	termTrigger := NewEmptyTermTrigger()

	navGroup := NewCommandGroup("Navigation and Control")
	navGroup.AddCommand(termbox.KeyCtrlC, SimpleFunc{
		Fn: func() {
			log.Printf("quitting")
			if termTrigger != nil {
				termTrigger.Stop()
			}
			if controls.OnQuit != nil {
				controls.OnQuit()
			}
			os.Exit(0)
		},
		Descr: "Quit the application",
	})
	navGroup.AddCommand(termbox.Key('l'), SimpleFunc{
		Fn: func() {
			*controls.LoggingEnabled = !*controls.LoggingEnabled
			if *controls.LoggingEnabled {
				log.Printf("=== BEAT LOGGING ENABLED")
			} else {
				log.Printf("=== BEAT LOGGING DISABLED")
			}
		},
		Descr: "Toggle beat logging on/off",
	})

	clockGroup := NewCommandGroup("Clock Control")
	clockGroup.AddCommand(termbox.Key('t'), SimpleFunc{
		Fn: func() {
			controls.FrequencyKeeper.AddTime(time.Now())
			bpm, _ := controls.FrequencyKeeper.GetBpm()
			log.Printf("=== TAP bpm: %f", bpm)
		},
		Descr: "Tap to set tempo",
	})
	clockGroup.AddCommand(termbox.Key('s'), SimpleFunc{
		Fn: func() {
			bpm, _ := controls.FrequencyKeeper.GetBpm()
			controls.Clock.SetBpm(bpm)
			log.Printf("=== SYNC bpm: %f", bpm)
		},
		Descr: "Sync clock to tapped tempo",
	})
	clockGroup.AddCommand(termbox.Key('r'), SimpleFunc{
		Fn: func() {
			controls.Clock.Reset()
			log.Printf("=== RESET")
		},
		Descr: "Reset the clock",
	})
	clockGroup.AddStringCommand(termbox.Key('m'), StringFunc{
		Fn: func(s string) {
			bpm, err := strconv.ParseFloat(s, 64)
			if err == nil {
				controls.Clock.SetBpm(bpm)
				log.Printf("=== SET bpm: %f", bpm)
			}
		},
		Descr: "Set BPM to specific value",
	})

	tempoGroup := NewCommandGroup("Tempo Adjustment")
	tempoGroup.AddCommand(termbox.Key('['), SimpleFunc{
		Fn: func() {
			controls.Clock.Nudge(-10 * time.Millisecond)
			log.Printf("=== NUDGE BACK -10ms")
		},
		Descr: "Nudge clock backward by 10ms",
	})
	tempoGroup.AddCommand(termbox.Key(']'), SimpleFunc{
		Fn: func() {
			controls.Clock.Nudge(10 * time.Millisecond)
			log.Printf("=== NUDGE FORWARD +10ms")
		},
		Descr: "Nudge clock forward by 10ms",
	})
	tempoGroup.AddCommand(termbox.Key('+'), SimpleFunc{
		Fn: func() {
			controls.Clock.Pace(.01)
			log.Printf("=== PACE UP .01")
		},
		Descr: "Increase tempo by 0.01 BPM",
	})
	tempoGroup.AddCommand(termbox.Key('-'), SimpleFunc{
		Fn: func() {
			controls.Clock.Pace(-.01)
			log.Printf("=== PACE DOWN .01")
		},
		Descr: "Decrease tempo by 0.01 BPM",
	})

	termTrigger.AddGroup(navGroup)
	termTrigger.AddGroup(clockGroup)
	termTrigger.AddGroup(tempoGroup)

	return termTrigger
}
