package controls

import (
	"encoding/json"
	"reflect"
)

// TODO: Make this a generic type with numerical dial as a parameter

type DialMap struct {
	Dials   map[string]*NumericDial
	channel chan ValueMap
}

func (m DialMap) SetValue(values ValueMap) {
	for name, value := range values {
		m.Dials[name].SetValue(value)
	}
}

func (m DialMap) GetString() string {
	r, err := json.Marshal(m)
	if err != nil {
		return err.Error()
	}
	return string(r)
}

func (m DialMap) GetValue() ValueMap {
	res := ValueMap{}
	for name, dial := range m.Dials {
		res[name] = dial.Value
	}
	return res
}

func (m DialMap) Emit() {
	m.channel <- m.GetValue()
}

func NewNumericDialMap(channels ...string) DialMap {
	res := DialMap{
		Dials:   make(map[string]*NumericDial),
		channel: make(chan ValueMap),
	}
	for _, channel := range channels {
		res.Dials[channel] = NewNumericDial()
	}

	go func() {
		for {
			cases := make([]reflect.SelectCase, 0, len(res.Dials))
			for _, dial := range res.Dials {
				cases = append(cases, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(dial.Channel()),
				})
			}
			_, _, ok := reflect.Select(cases)
			if ok {
				res.Emit()
			}
		}
	}()

	return res
}

func (m DialMap) Channel() <-chan ValueMap {
	return m.channel
}

func (m DialMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Dials)
}
