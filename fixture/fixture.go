package fixture

import (
	"github.com/mholzen/play-go/controls"
)

type Fixture interface {
	GetValues() []byte
	GetChannels() []string

	GetChannelValue(channel string) byte
	SetChannelValue(channel string, value byte)

	GetChannelValues() controls.ChannelValues
	SetChannelValues(values controls.ChannelValues)

	SetAll(value byte)
}

type ChannelFixture struct {
	Model  *ModelChannels
	Values []byte
}

func (f ChannelFixture) setChannelValue(channel string, value byte) {
	i, err := f.Model.GetAddress(channel)
	if err != nil {
		// log.Printf("cannot set value for '%s': %s", f.Model.Name, err)
		return
	}
	f.Values[i] = value
}

func (f ChannelFixture) GetChannelValue(channel string) byte {
	i, err := f.Model.GetAddress(channel)
	if err != nil {
		return 0 // TODO: change the method signature
	}
	return f.Values[i]
}

func (f ChannelFixture) SetChannelValue(channel string, value byte) {
	f.setChannelValue(channel, value)
}

func (f ChannelFixture) SetChannelValues(values controls.ChannelValues) {
	for channel, value := range values {
		f.SetChannelValue(channel, value)
	}
}

func (f ChannelFixture) GetChannelValues() controls.ChannelValues {
	res := make(controls.ChannelValues)
	for channel, index := range f.Model.IndexByChannel {
		res[channel] = f.Values[index]
	}
	return res
}

func (f ChannelFixture) SetAll(value byte) {
	for _, channel := range f.Model.Channels {
		f.SetChannelValue(channel, value)
	}
}

func (f ChannelFixture) GetChannels() []string {
	return f.Model.Channels
}

func (f ChannelFixture) GetValues() []byte {
	return f.Values
}

func (f ChannelFixture) Clone() ChannelFixture {
	return ChannelFixture{
		Model:  f.Model,
		Values: append([]byte{}, f.Values...),
	}
}

func NewChannelFixture(model ModelChannels) ChannelFixture {
	return ChannelFixture{
		Model:  &model,
		Values: make([]byte, len(model.Channels)),
	}
}

func ApplyTo(values controls.ChannelValues, f Fixture) {
	for k, v := range values {
		f.SetChannelValue(k, v)
	}
}
