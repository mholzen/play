package controls

import (
	"encoding/json"
	"fmt"
)

type ObservableNumericDialMap struct {
	Dials *NumericDialMap
	Observers[ChannelValues]
}

func NewObservableNumericDialMap(channels ...string) *ObservableNumericDialMap {
	return &ObservableNumericDialMap{
		Observers: *NewObservable[ChannelValues](),
		Dials:     NewNumericDialMap(channels...),
	}
}

func (m *ObservableNumericDialMap) GetItem(name string) (Item, error) {
	return m.Dials.GetItem(name)
}

func (m *ObservableNumericDialMap) Items() map[string]Item {
	return m.Dials.Items()
}

func (m *ObservableNumericDialMap) SetValue(values ChannelValues) {
	for name, value := range values {
		if dial, ok := (*m.Dials)[name]; ok {
			dial.SetValue(value)
		} else {
			panic(fmt.Sprintf("dial '%s' not found", name))
		}
	}
	m.Notify(values)
}

func (m *ObservableNumericDialMap) SetChannelValue(name string, value byte) {
	if dial, ok := (*m.Dials)[name]; ok {
		dial.SetValue(value)
	} else {
		panic(fmt.Sprintf("dial '%s' not found", name))
	}
	m.Notify(m.GetValue())
}

func (m *ObservableNumericDialMap) GetValue() ChannelValues {
	res := ChannelValues{}
	for name, dial := range *m.Dials {
		res[name] = dial.Value
	}
	return res
}

func (m *ObservableNumericDialMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Dials)
}
