package fixture

import (
	"fmt"
	"log"
	"time"

	"github.com/bep/debounce"
)

type FixtureI interface {
	SetValue(channel string, value byte)
	SetAll(value byte)
	GetValues() []byte
	SetOnUpdate(func(FixtureI), time.Duration)
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

type Fixture struct {
	Model    ModelChannels
	Values   []byte
	onUpdate func()
}

func (f *Fixture) SetOnUpdate(onUpdate func(FixtureI), debounceTime time.Duration) {
	onUpdateThis := func() {
		onUpdate(f)
	}
	if debounceTime > 0 {
		d := debounce.New(debounceTime)
		f.onUpdate = func() {
			log.Printf("pre debound called")
			d(onUpdateThis)
		}
	} else {
		f.onUpdate = onUpdateThis
	}
}

func (f *Fixture) SetValue(channel string, value byte) {
	i, err := f.Model.GetAddress(channel)
	if err != nil {
		// log.Printf("cannot set value for '%s': %s", f.Model.Name, err)
		return
	}
	// log.Printf("setting '%s'[%s] to %d", f.Model.Name, channel, value)
	f.Values[i] = value

	if f.onUpdate != nil {
		f.onUpdate()
	}
}
func (f *Fixture) SetAll(value byte) {
	for _, channel := range f.Model.Channels {
		f.SetValue(channel, value)
	}
}

func (f *Fixture) GetValues() []byte {
	return f.Values
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
