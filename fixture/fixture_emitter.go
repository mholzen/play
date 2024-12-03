package fixture

import (
	"github.com/mholzen/play-go/controls"
)

type FixturesEmitter struct {
	Fixtures Fixtures
	channel  chan FixtureValues
}

func NewFixtureEmitter(fixtures Fixtures) FixturesEmitter {
	return FixturesEmitter{Fixtures: fixtures}
}

func (f *FixturesEmitter) Channel() <-chan FixtureValues {
	if f.channel == nil {
		// log.Printf("creating channel")
		f.channel = make(chan FixtureValues)
	}
	return f.channel
}

func (f *FixturesEmitter) GetValue() FixtureValues {
	return f.Fixtures.GetValue()
}

func (f *FixturesEmitter) SetChannelValue(channel string, value byte) {
	f.Fixtures.SetChannelValue(channel, value)
	f.Emit()
}

func (f *FixturesEmitter) SetAll(value byte) {
	f.Fixtures.SetAll(value)
	f.Emit()
}

func (f *FixturesEmitter) SetValueMap(values controls.ValueMap) {
	for channel, value := range values {
		f.Fixtures.SetChannelValue(channel, value)
	}
	f.Emit()
}

func (f *FixturesEmitter) GetValues() []byte {
	return f.Fixtures.GetValues()
}

func (f *FixturesEmitter) Emit() {
	select {
	case f.channel <- f.GetValue():
		// log.Printf("emitter %v emitting %v", f.Fixture, f.GetValueMap())
	default:
		// log.Printf("no receiver")
	}
}

func (f *FixturesEmitter) GetValueMap() controls.ValueMap {
	return f.Fixtures.GetValueMap()
}
