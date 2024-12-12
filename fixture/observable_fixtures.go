package fixture

import (
	"github.com/mholzen/play-go/controls"
)

type ObservableFixtures struct {
	controls.Observable[FixtureValues]
	Fixtures FixturesGeneric[FixtureI]
}

func (f *ObservableFixtures) SetValueMap(values controls.ValueMap) {
	for _, fixture := range f.Fixtures {
		fixture.SetValueMap(values)
	}
	f.Notify(f.GetFixtureValues())
}

func (f *ObservableFixtures) GetChannels() []string {
	panic("not implemented")
}

func (f *ObservableFixtures) GetFixtureValues() FixtureValues {
	values := make(FixtureValues)
	for address, fixture := range f.Fixtures {
		values[address] = fixture.GetValueMap()
	}
	return values
}

func (f *ObservableFixtures) GetFixtures() FixturesGeneric[FixtureI] {
	return f.Fixtures
}
