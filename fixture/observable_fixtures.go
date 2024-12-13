package fixture

import (
	"log"

	"github.com/mholzen/play-go/controls"
)

type ObservableFixtures struct {
	FixturesGeneric[FixtureI]
	controls.Observable[FixtureValues]
}

func NewObservableFixtures(fixtures FixturesInterface[FixtureI]) *ObservableFixtures {
	res := &ObservableFixtures{
		FixturesGeneric: *NewFixturesGeneric[FixtureI](),
		Observable:      *controls.NewObservable[FixtureValues](),
	}
	for address, fixture := range fixtures.GetFixtures() {
		res.AddFixture(fixture, address)
	}
	return res
}

func NewIndividualObservableFixtures(fixtures FixturesInterface[FixtureI]) *ObservableFixtures {
	res := &ObservableFixtures{
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

func (f *ObservableFixtures) SetValueMap(values controls.ValueMap) {
	for _, fixture := range f.FixturesGeneric {
		fixture.SetValueMap(values)
	}
	f.Notify(f.FixturesGeneric.GetValue())
}

func NewObservableDialMapForAllChannels(channels []string, fixtures *ObservableFixtures) *controls.ObservableDialMap {
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
