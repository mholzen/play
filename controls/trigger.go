package controls

type TriggerFunc func(Clock) bool

func TriggerOnTick(tick int) TriggerFunc {
	return func(c Clock) bool {
		return c.Tick() == tick
	}
}

func TriggerOnBeat(beat int) TriggerFunc {
	return func(c Clock) bool {
		return c.Beat() == beat && c.Tick() == 0
	}
}

func TriggerOnBeats() TriggerFunc {
	return TriggerOnTick(0)
}

func TriggerOnBar(bar int) TriggerFunc {
	return func(c Clock) bool {
		return c.Bar() == bar && c.Beat() == 0 && c.Tick() == 0
	}
}

func TriggerOnBars() TriggerFunc {
	return TriggerOnBar(0)
}

func TriggerOnPhrase(phrase int) TriggerFunc {
	return func(c Clock) bool {
		return c.Phrase() == phrase && c.Bar() == 0 && c.Beat() == 0 && c.Tick() == 0
	}
}

func TriggerOnPhrases() TriggerFunc {
	return TriggerOnPhrase(0)
}

type Trigger struct {
	When    TriggerFunc
	Enabled bool
	Do      func()
}

type Triggers []Trigger

func (t Triggers) Enable() {
	for _, trigger := range t {
		trigger.Enabled = true
	}
}

func (t Triggers) Disable() {
	for _, trigger := range t {
		trigger.Enabled = false
	}
}
