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

func TriggerOnPhraseRatio(numerator int, denominator int) TriggerFunc {
	return func(c Clock) bool {
		ticks := numerator * (c.TicksPerBeat * c.BeatsPerBar * c.BarPerPhrase) / denominator
		trigger := c.Ticks()%ticks == 0
		return trigger
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
	When    TriggerFunc `json:"-"`
	Enabled bool        `json:"enabled"`
	Do      func()      `json:"-"`
}

type Triggers []Trigger

func (t Triggers) Enable() {
	for i := range t {
		t[i].Enabled = true
	}
}

func (t Triggers) Disable() {
	for i := range t {
		t[i].Enabled = false
	}
}
