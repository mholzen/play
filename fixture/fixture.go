package fixture

import (
	"github.com/mholzen/play-go/controls"
)

type FixtureI interface {
	SetChannelValue(channel string, value byte)
	GetValueMap() controls.ValueMap // TODO: how would FixturesInterface implement this?
	SetValueMap(values controls.ValueMap)
	SetAll(value byte)
	GetValues() []byte
	GetChannels() []string
	// Clone() FixtureI
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
	// log.Printf("setting '%s'[%s] to %d", f.Model.Name, channel, value)
	f.Values[i] = value
}

func (f Fixture) SetChannelValue(channel string, value byte) {
	f.setChannelValue(channel, value)
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

func (f Fixture) GetValueMap() controls.ValueMap {
	res := make(controls.ValueMap)
	for channel, index := range f.Model.IndexByChannel {
		res[channel] = f.Values[index]
	}
	return res
}

func (f Fixture) SetValueMap(values controls.ValueMap) {
	for channel, value := range values {
		f.SetChannelValue(channel, value)
	}
}

func ApplyTo(values controls.ValueMap, f FixtureI) {
	for k, v := range values {
		f.SetChannelValue(k, v)
	}
}

func (f Fixture) Clone() Fixture {
	return Fixture{
		Model:  f.Model,
		Values: append([]byte{}, f.Values...),
	}
}
