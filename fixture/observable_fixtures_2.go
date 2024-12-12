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
		var fixtureI FixtureI = observableFixture
		res.AddFixture(fixtureI, address)
		observableFixture.AddObserver(ch)
	}
	go func() {
		for range ch {
			res.Notify(res.GetValue())
		}
	}()
	return res
}

func (f *ObservableFixtures2) SetValueMap(values controls.ValueMap) {
	for _, fixture := range f.FixturesGeneric {
		// log.Printf("ObservableFixtures2: setting fixture %d with value %v", address, values)
		fixture.SetValueMap(values)
	}
	// log.Printf("ObservableFixtures2: notifying observers of %v", f.FixturesGeneric.GetValue())
	f.Notify(f.FixturesGeneric.GetValue())
}

func NewObservableDialMapForAllChannels(channels []string, fixtures *ObservableFixtures2) *controls.ObservableDialMap {
	dialMap := controls.NewObservableNumericDialMap(channels...)
	received := make(chan controls.ValueMap)
	dialMap.AddObserver(received)
	go func() {
		for valueMap := range received {
			log.Printf("dial map received value map %v", valueMap)
			fixtures.SetValueMap(valueMap)
		}
	}()
	return dialMap
}
