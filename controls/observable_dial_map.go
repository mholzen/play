package controls

import (
	"encoding/json"
	"fmt"
)

type ObservableDialMap struct {
	Dials *ObservableNumericDialmap
	Observable[ChannelValues]
}

func (m *ObservableDialMap) GetItem(name string) (Item, error) {
	return m.Dials.GetItem(name)
}

func (m *ObservableDialMap) Items() map[string]Item {
	return m.Dials.Items()
}

func (m *ObservableDialMap) SetValue(values ChannelValues) {
	m.Dials.SetValue(values)
}

type ObservableNumericDialMap struct {
	Observers[ChannelValues]
	Dials *NumericDialMap
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

func _NewObservableNumericDialMap(channels ...string) *ObservableDialMap {
	dialMap := NewObservableNumericDialMap2(channels...)

	ch := make(chan byte)
	for _, dial := range *dialMap {
		dial.AddObserver(ch)
	}

	res := &ObservableDialMap{
		Observable: NewObservable[ChannelValues](),
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

func (m *ObservableNumericDialMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Dials)
}

type ObservableDialMap2 struct { // TODO: reconcile with ObservableDialMap
	Observers[ChannelValues]
	ItemMap // TODO: replace with DialMap
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

func NewObservableDialMap2() *ObservableDialMap2 {
	res := &ObservableDialMap2{
		Observers: *NewObservable[ChannelValues](),
		ItemMap:   NewItemMap(),
	}
	return res
}
