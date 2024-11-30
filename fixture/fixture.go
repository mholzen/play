package fixture

import (
	"fmt"
	"log"

	"github.com/mholzen/play-go/controls"
)

type FixtureI interface {
	SetChannelValue(channel string, value byte)
	SetAll(value byte)
	GetValues() []byte
}

type ModelChannels struct {
	Name           string
	Channels       []string
	IndexByChannel map[string]int
}

func NewModelChannels(name string, channels []string) ModelChannels {
	m := ModelChannels{Name: name}
	m.SetChannels(channels)
	return m
}

type Fixture struct { // TODO: merge this with InstalledFixture
	Model   ModelChannels
	Values  []byte
	channel chan controls.ValueMap
}

func (f *Fixture) setChannelValue(channel string, value byte) {
	i, err := f.Model.GetAddress(channel)
	if err != nil {
		// log.Printf("cannot set value for '%s': %s", f.Model.Name, err)
		return
	}
	// log.Printf("setting '%s'[%s] to %d", f.Model.Name, channel, value)
	f.Values[i] = value
}

func (f *Fixture) Emit() {
	select {
	case f.channel <- f.GetValueMap():
		log.Printf("emitting %v", f.GetValueMap())
	default:
	}
}

func (f *Fixture) SetChannelValue(channel string, value byte) {
	f.setChannelValue(channel, value)
	f.Emit()
}

func (f *Fixture) SetAll(value byte) {
	for _, channel := range f.Model.Channels {
		f.SetChannelValue(channel, value)
	}
	f.Emit()
}

func (f *Fixture) GetValues() []byte {
	return f.Values
}

func (f *Fixture) GetValueMap() controls.ValueMap {
	res := make(controls.ValueMap)
	for channel, index := range f.Model.IndexByChannel {
		res[channel] = f.Values[index]
	}
	return res
}

func (f *Fixture) Channel() chan controls.ValueMap {
	if f.channel == nil {
		f.channel = make(chan controls.ValueMap)
	}
	return f.channel
}

func (m ModelChannels) GetAddress(name string) (int, error) {
	if res, ok := m.IndexByChannel[name]; ok {
		return res, nil
	} else {
		return 0, fmt.Errorf("cannot find channel '%s'", name)
	}
}

func (m *ModelChannels) SetChannels(channels []string) {
	m.Channels = channels
	m.IndexByChannel = ArrayToIndex(m.Channels)
}

func ApplyTo(values controls.ValueMap, f FixtureI) {
	for k, v := range values {
		f.SetChannelValue(k, v)
	}
}
