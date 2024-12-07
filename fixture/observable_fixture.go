package fixture

import "github.com/mholzen/play-go/controls"

type ObservableFixture struct {
	controls.Observable[controls.ValueMap]
	Fixture FixtureI
}

func NewObservableFixture(initial FixtureI) *ObservableFixture {
	return &ObservableFixture{
		Fixture:    initial,
		Observable: *controls.NewObservable[controls.ValueMap](),
	}
}

func (f *ObservableFixture) GetValueMap() controls.ValueMap {
	return f.Fixture.GetValueMap()
}

func (f *ObservableFixture) SetValueMap(values controls.ValueMap) {
	f.Fixture.SetValueMap(values)
	f.Notify(values)
}

func (f *ObservableFixture) GetChannels() []string {
	return f.Fixture.GetChannels()
}

func (f *ObservableFixture) GetValues() []byte {
	return f.Fixture.GetValues()
}

func (f *ObservableFixture) SetAll(value byte) {
	f.Fixture.SetAll(value)
	f.Notify(f.GetValueMap())
}

func (f *ObservableFixture) SetChannelValue(name string, value byte) {
	f.Fixture.SetChannelValue(name, value)
	f.Notify(f.GetValueMap())
}
