package fixture

import (
	"log"

	"github.com/mholzen/play-go/controls"
)

type ObservableFixtures struct {
	AddressableFixtures[Fixture]
	controls.Observers[FixtureValues]
}

func NewObservableFixtures(fixtures Fixtures[Fixture]) *ObservableFixtures {
	res := &ObservableFixtures{
		AddressableFixtures: *NewAddressableFixtures[Fixture](),
		Observers:           *controls.NewObservable[FixtureValues](),
	}
	for address, fixture := range fixtures.GetFixtures() {
		res.AddFixture(fixture, address)
	}
	return res
}

func NewIndividualObservableFixtures(fixtures Fixtures[Fixture]) *ObservableFixtures {
	res := &ObservableFixtures{
		AddressableFixtures: *NewAddressableFixtures[Fixture](),
		Observers:           *controls.NewObservable[FixtureValues](),
	}

	ch := make(chan controls.ChannelValues)
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

func (f *ObservableFixtures) SetChannelValues(values controls.ChannelValues) {
	for _, fixture := range f.AddressableFixtures {
		fixture.SetChannelValues(values)
	}
	f.Notify(f.AddressableFixtures.GetValue())
}

func NewObservableDialMapForAllChannels(fixtures *ObservableFixtures) *controls.ObservableDialMap {
	channels := fixtures.GetChannels()
	dialMap := controls.NewObservableNumericDialMap(channels...)
	received := make(chan controls.ChannelValues)
	dialMap.AddObserver(received)
	go func() {
		for valueMap := range received {
			log.Printf("dial map received value map %v", valueMap)
			fixtures.SetChannelValues(valueMap)
		}
	}()
	return dialMap
}
