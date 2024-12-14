package controls

import (
	"encoding/json"
	"fmt"
)

// TODO: Make this a generic type with numerical dial as a parameter

type ObservableDialMap struct {
	Observable[ChannelValues]
	Dials *NumericDialMap
}

func (m *ObservableDialMap) SetValue(values ChannelValues) {
	for name, value := range values {
		if dial, ok := (*m.Dials)[name]; ok {
			dial.SetValue(value)
		} else {
			panic(fmt.Sprintf("dial '%s' not found", name))
		}
	}
	m.Notify(values)
}

func (m *ObservableDialMap) SetChannelValue(name string, value byte) {
	if dial, ok := (*m.Dials)[name]; ok {
		dial.SetValue(value)
	} else {
		panic(fmt.Sprintf("dial '%s' not found", name))
	}
	m.Notify(m.GetValue())
}

func (m *ObservableDialMap) GetValue() ChannelValues {
	res := ChannelValues{}
	for name, dial := range *m.Dials {
		res[name] = dial.Value
	}
	return res
}

func NewObservableNumericDialMap(channels ...string) *ObservableDialMap {
	return &ObservableDialMap{
		Observable: *NewObservable[ChannelValues](),
		Dials:      NewNumericDialMap2(channels...),
	}
}

func (m *ObservableDialMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Dials)
}
