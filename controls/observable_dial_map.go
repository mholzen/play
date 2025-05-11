package controls

import (
	"encoding/json"
	"fmt"
)

// TODO: Make this a generic type with numerical dial as a parameter

type ObservableDialMap struct {
	Observable[ChannelValues]
	Dials *ObservableNumericalDialmap
}

func (m *ObservableDialMap) GetItem(name string) (Item, error) {
	return m.Dials.GetItem(name)
}

func (m *ObservableDialMap) Items() map[string]Item {
	return m.Dials.Items()
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
	dialMap := NewObservableNumericDialMap2(channels...)

	ch := make(chan byte)
	for _, dial := range *dialMap {
		dial.AddObserver(ch)
	}

	res := &ObservableDialMap{
		Observable: *NewObservable[ChannelValues](),
		Dials:      dialMap,
	}

	go func() {
		for range ch {
			channelValues := dialMap.GetChannelValues()
			res.Observable.Notify(channelValues)
		}
	}()

	return res
}

func (m *ObservableDialMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Dials)
}
