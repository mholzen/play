package controls

import (
	"encoding/json"
	"fmt"
)

type ObservableNumericalDialMap struct {
	Observers[ChannelValues]
	Dials *NumericDialMap
}

func (m *ObservableNumericalDialMap) GetItem(name string) (Item, error) {
	return m.Dials.GetItem(name)
}

func (m *ObservableNumericalDialMap) Items() map[string]Item {
	return m.Dials.Items()
}

func (m *ObservableNumericalDialMap) SetValue(values ChannelValues) {
	for name, value := range values {
		if dial, ok := (*m.Dials)[name]; ok {
			dial.SetValue(value)
		} else {
			panic(fmt.Sprintf("dial '%s' not found", name))
		}
	}
	m.Notify(values)
}

func (m *ObservableNumericalDialMap) SetChannelValue(name string, value byte) {
	if dial, ok := (*m.Dials)[name]; ok {
		dial.SetValue(value)
	} else {
		panic(fmt.Sprintf("dial '%s' not found", name))
	}
	m.Notify(m.GetValue())
}

func (m *ObservableNumericalDialMap) GetValue() ChannelValues {
	res := ChannelValues{}
	for name, dial := range *m.Dials {
		res[name] = dial.Value
	}
	return res
}

func (m *ObservableNumericalDialMap) GetChannels() []string {
	return m.Dials.GetChannels()
}

// only used in tests
func NewObservableNumericDialMap(channels ...string) *ObservableNumericalDialMap {
	return &ObservableNumericalDialMap{
		Observers: *NewObservable[ChannelValues](),
		Dials:     NewNumericDialMap(channels...),
	}
}

func (m *ObservableNumericalDialMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Dials)
}

// Used in NewObservableDialMapForAllChannels
type ObservableDialMap2 struct {
	Observers[ChannelValues]
	ItemMap // TODO: replace with DialMap
	// DialMap
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

func NewObservableDialMap2(channels ...string) *ObservableDialMap2 {
	res := &ObservableDialMap2{
		Observers: *NewObservable[ChannelValues](),
		ItemMap:   NewItemMap(),
	}
	for _, channel := range channels {
		res.AddItem(channel, NewObservableNumericalDial(NewNumericDial()))
	}
	return res
}
