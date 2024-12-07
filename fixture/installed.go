package fixture

import (
	"github.com/mholzen/play-go/controls"
)

type InstalledFixture struct {
	Fixture FixtureI
	Address int
}

func (f *InstalledFixture) GetValues() []byte {
	return f.Fixture.GetValues()
}

func (f *InstalledFixture) SetChannelValue(name string, value byte) {
	f.Fixture.SetChannelValue(name, value)
}

func (f *InstalledFixture) SetAll(value byte) {
	f.Fixture.SetAll(value)
}

func (f *InstalledFixture) SetValues(values controls.ValueMap) {
	for channel, value := range values {
		f.SetChannelValue(channel, value)
	}
}

func (f *InstalledFixture) SetValueMap(values controls.ValueMap) {
	f.SetValues(values)
}

func (f *InstalledFixture) GetValueMap() controls.ValueMap {
	return f.Fixture.GetValueMap()
}

func (f *InstalledFixture) GetChannels() []string {
	return f.Fixture.GetChannels()
}
