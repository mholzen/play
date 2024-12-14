package fixture

import (
	"github.com/mholzen/play-go/controls"
)

type FixtureI interface {
	SetChannelValue(channel string, value byte)
	GetChannelValues() controls.ChannelValues
	SetChannelValues(values controls.ChannelValues)
	SetAll(value byte)
	GetValues() []byte
	GetChannels() []string
}

type Fixture struct {
	Model  *ModelChannels
	Values []byte
}

func (f Fixture) setChannelValue(channel string, value byte) {
	i, err := f.Model.GetAddress(channel)
	if err != nil {
		// log.Printf("cannot set value for '%s': %s", f.Model.Name, err)
		return
	}
	f.Values[i] = value
}

func (f Fixture) SetChannelValue(channel string, value byte) {
	f.setChannelValue(channel, value)
}

func (f Fixture) SetChannelValues(values controls.ChannelValues) {
	for channel, value := range values {
		f.SetChannelValue(channel, value)
	}
}

func (f Fixture) GetChannelValues() controls.ChannelValues {
	res := make(controls.ChannelValues)
	for channel, index := range f.Model.IndexByChannel {
		res[channel] = f.Values[index]
	}
	return res
}

func (f Fixture) SetAll(value byte) {
	for _, channel := range f.Model.Channels {
		f.SetChannelValue(channel, value)
	}
}

func (f Fixture) GetChannels() []string {
	return f.Model.Channels
}

func (f Fixture) GetValues() []byte {
	return f.Values
}

func (f Fixture) Clone() Fixture {
	return Fixture{
		Model:  f.Model,
		Values: append([]byte{}, f.Values...),
	}
}

func ApplyTo(values controls.ChannelValues, f FixtureI) {
	for k, v := range values {
		f.SetChannelValue(k, v)
	}
}
