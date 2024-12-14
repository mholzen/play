package controls

import (
	"encoding/json"
	"fmt"
)

// TODO: Make this a generic type with numerical dial as a parameter

type DialMap struct {
	Dials   map[string]*NumericDial
	channel chan ChannelValues
}

func (m DialMap) SetValue(values ChannelValues) {
	for name, value := range values {
		if dial, ok := m.Dials[name]; ok {
			dial.SetValue(value)
		} else {
			panic(fmt.Sprintf("dial '%s' not found", name))
		}
	}
}

func (m DialMap) SetChannelValue(name string, value byte) {
	if dial, ok := m.Dials[name]; ok {
		dial.SetValue(value)
	} else {
		panic(fmt.Sprintf("dial '%s' not found", name))
	}
}

func (m DialMap) GetString() string {
	r, err := json.Marshal(m)
	if err != nil {
		return err.Error()
	}
	return string(r)
}

func (m DialMap) GetValue() ChannelValues {
	res := ChannelValues{}
	for name, dial := range m.Dials {
		res[name] = dial.Value
	}
	return res
}

// func (m DialMap) Emit() {
// 	m.channel <- m.GetValue()
// }

func NewNumericDialMap(channels ...string) *DialMap {
	res := DialMap{
		Dials:   make(map[string]*NumericDial),
		channel: make(chan ChannelValues),
	}
	for _, channel := range channels {
		res.Dials[channel] = NewNumericDial()
	}

	// valueMap := make(ValueMap)
	// for channelName, dial := range res.Dials {
	// 	go func(channelName string, dial *NumericDial) {
	// 		for value := range dial.Channel() {
	// 			valueMap[channelName] = value
	// 			res.channel <- valueMap
	// 		}
	// 	}(channelName, dial)
	// }

	return &res
}

// func (m DialMap) Channel() <-chan ValueMap {
// 	return m.channel
// }

func (m DialMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Dials)
}

func NewDialMap() *DialMap {
	return NewNumericDialMap("mode", "dimmer", "strobe", "tilt", "pan", "speed", "r", "g", "b", "w", "a", "uv")
}
