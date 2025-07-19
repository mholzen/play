package fixture

import (
	"github.com/mholzen/play/controls"
)

type ObservableFixtures struct {
	AddressableFixtures[Fixture]
	controls.Observers[FixtureValues]
}

func newObservableFixtures() *ObservableFixtures {
	return &ObservableFixtures{
		AddressableFixtures: *NewAddressableFixtures[Fixture](),
		Observers:           *controls.NewObservable[FixtureValues](),
	}
}

func NewObservableFixtures(fixtures Fixtures[Fixture]) *ObservableFixtures {
	res := newObservableFixtures()
	for address, fixture := range fixtures.GetFixtures() {
		res.AddFixture(fixture, address)
	}
	return res
}

func NewIndividualObservableFixtures(fixtures Fixtures[Fixture]) *ObservableFixtures {
	res := newObservableFixtures()

	ch := make(chan controls.ChannelValues)
	for address, fixture := range fixtures.GetFixtures() {
		observableFixture := NewObservableFixture(fixture)
		observableFixture.AddObserver(ch)

		res.AddFixture(observableFixture, address)
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

func (f *ObservableFixtures) SetChannelValue(channel string, value byte) {
	f.AddressableFixtures.SetChannelValue(channel, value)
	f.Notify(f.AddressableFixtures.GetValue())
}
