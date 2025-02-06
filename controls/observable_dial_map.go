package controls

import (
	"encoding/json"
	"fmt"
)

type ObservableDialMap struct {
	Observers[ChannelValues]
	Dials *NumericDialMap
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

func (m *ObservableDialMap) GetChannels() []string {
	return m.Dials.GetChannels()
}

func NewObservableNumericDialMap(channels ...string) *ObservableDialMap {
	return &ObservableDialMap{
		Observers: *NewObservable[ChannelValues](),
		Dials:     NewNumericDialMap(channels...),
	}
}

func (m *ObservableDialMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Dials)
}

// NEW
type ObservableDialMap2 struct {
	Observers[ChannelValues]
	ItemMap
}

func (m *ObservableDialMap2) AddItem(name string, item Item) {
	m.ItemMap[name] = item

	ch := make(chan ChannelValues)
	if observable, ok := item.(Observable[ChannelValues]); ok {
		observable.AddObserver(ch)
	}
	go func() {
		for values := range ch {
			m.Notify(values)
		}
	}()

	ch2 := make(chan byte)
	if observable, ok := item.(Observable[byte]); ok {
		observable.AddObserver(ch2)
	}
	go func() {
		for value := range ch2 {
			values := ChannelValues{}
			values[name] = value
			m.Notify(values)
		}
	}()

}

func (m *ObservableDialMap2) SetChannelValue(name string, value byte) {
	if item, ok := m.ItemMap[name]; ok {
		if observable, ok := item.(Observable[byte]); ok {
			if dial, ok := observable.(Settable); ok {
				dial.SetValue(value)
			}
		}
	}
}

func NewObservableDialMap2() *ObservableDialMap2 {
	return &ObservableDialMap2{
		Observers: *NewObservable[ChannelValues](),
		ItemMap:   NewItemMap(),
	}
}
