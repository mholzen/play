package fixture

import (
	"github.com/mholzen/play-go/controls"
)

type ObservableFixtures struct {
	controls.Observable[FixtureValues]
	Fixtures FixturesGeneric[FixtureI]
}

func NewObservableFixtures(fixtures Fixtures) *ObservableFixtures {
	f := &ObservableFixtures{
		Fixtures:   *NewFixturesGeneric[FixtureI](),
		Observable: *controls.NewObservable[FixtureValues](),
	}
	ch := make(chan controls.ValueMap)
	for _, fixture := range fixtures {
		observableFixture := NewObservableFixture(fixture.Fixture)
		f.Fixtures.AddFixture(observableFixture, fixture.Address)
		observableFixture.AddObserver(ch)
	}
	go func() {
		for range ch {
			f.Notify(f.GetFixtureValues())
		}
	}()
	return f
}

// func (f *ObservableFixtures) SetValue(values FixtureValues) {
// 	for _, valueMap := range values {
// 		for _, observable := range f.Fixtures {
// 			observable.Fixture.SetValueMap(values.ValueMap)
// 		}
// 	}
// 	f.Notify(values)
// }

func (f *ObservableFixtures) SetValueMap(values controls.ValueMap) {
	for _, fixture := range f.Fixtures {
		(*fixture).SetValueMap(values)
	}
	f.Notify(f.GetFixtureValues())
}

func (f *ObservableFixtures) GetChannels() []string {
	panic("not implemented")
}

func (f *ObservableFixtures) GetFixtureValues() FixtureValues {
	values := make(FixtureValues)
	for address, fixture := range f.Fixtures {
		values[address] = (*fixture).GetValueMap()
	}
	return values
}

func (f *ObservableFixtures) GetFixtures() FixturesGeneric[FixtureI] {
	return f.Fixtures
}
