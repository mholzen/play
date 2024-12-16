package fixture

import "github.com/mholzen/play-go/controls"

type ObservableFixture struct {
	controls.Observers[controls.ChannelValues]
	Fixture Fixture
}

func NewObservableFixture(initial Fixture) *ObservableFixture {
	return &ObservableFixture{
		Fixture:   initial,
		Observers: *controls.NewObservable[controls.ChannelValues](),
	}
}

func (f *ObservableFixture) GetChannelValues() controls.ChannelValues {
	return f.Fixture.GetChannelValues()
}

func (f *ObservableFixture) SetChannelValues(values controls.ChannelValues) {
	f.Fixture.SetChannelValues(values)
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
	f.Notify(f.GetChannelValues())
}

func (f *ObservableFixture) SetChannelValue(name string, value byte) {
	f.Fixture.SetChannelValue(name, value)
	f.Notify(f.GetChannelValues())
}
