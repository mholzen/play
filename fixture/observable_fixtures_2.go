package fixture

import (
	"log"

	"github.com/mholzen/play-go/controls"
)

type ObservableFixtures2 struct {
	FixturesGeneric[FixtureI]
	controls.Observable[FixtureValues]
}

func NewObservableFixtures2(fixtures FixturesInterface[FixtureI]) *ObservableFixtures2 {
	res := &ObservableFixtures2{
		FixturesGeneric: *NewFixturesGeneric[FixtureI](),
		Observable:      *controls.NewObservable[FixtureValues](),
	}
	for address, fixture := range fixtures.GetFixtures() {
		res.AddFixture(fixture, address)
	}
	return res
}

func NewIndividualObservableFixtures2(fixtures FixturesInterface[FixtureI]) *ObservableFixtures2 {
	res := &ObservableFixtures2{
		FixturesGeneric: *NewFixturesGeneric[FixtureI](),
		Observable:      *controls.NewObservable[FixtureValues](),
	}

	ch := make(chan controls.ValueMap)
	for address, fixture := range fixtures.GetFixtures() {
		observableFixture := NewObservableFixture(fixture)
		res.AddFixture(observableFixture, address)
		observableFixture.AddObserver(ch)
	}
	go func() {
		for range ch {
			res.Notify(res.GetValue())
		}
	}()
	return res
}

// func (f *ObservableFixtures) SetValue(values FixtureValues) {
// 	for _, valueMap := range values {
// 		for _, observable := range f.Fixtures {
// 			observable.Fixture.SetValueMap(values.ValueMap)
// 		}
// 	}
// 	f.Notify(values)
// }

func (f *ObservableFixtures2) SetValueMap(values controls.ValueMap) {
	for address, fixture := range f.FixturesGeneric {
		log.Printf("ObservableFixtures2: setting fixture %d with value %v", address, values)
		(*fixture).SetValueMap(values)
	}
	log.Printf("ObservableFixtures2: notifying observers of %v", f.FixturesGeneric.GetValue())
	f.Notify(f.FixturesGeneric.GetValue())
}

// func (f *ObservableFixtures) GetChannels() []string {
// 	panic("not implemented")
// }

// func (f *ObservableFixtures) GetFixtureValues() FixtureValues {
// 	values := make(FixtureValues)
// 	for address, fixture := range f.Fixtures {
// 		values[address] = (*fixture).GetValueMap()
// 	}
// 	return values
// }

// func (f *ObservableFixtures) GetFixtures() FixturesGeneric[FixtureI] {
// 	return f.Fixtures
// }
