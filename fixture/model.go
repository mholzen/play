package fixture

import (
	"fmt"

	"github.com/mholzen/play-go/controls"
)

type ModelChannels struct {
	Name           string
	Channels       []string
	IndexByChannel map[string]int
}

func NewModelChannels(channels ...string) ModelChannels {
	m := ModelChannels{}
	m.SetChannels(channels)
	return m
}

func NewModelChannelsWithName(name string, channels []string) ModelChannels {
	m := ModelChannels{Name: name}
	m.SetChannels(channels)
	return m
}

func (m ModelChannels) GetChannels() []string {
	return m.Channels
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

func (m ModelChannels) GetChannelValues(values []byte) controls.ChannelValues {
	res := make(controls.ChannelValues)
	for i, channel := range m.Channels {
		res[channel] = values[i]
	}
	return res
}

func (m ModelChannels) GetEmptyValues() []byte {
	return make([]byte, len(m.Channels))
}
