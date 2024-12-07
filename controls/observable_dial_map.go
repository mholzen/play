package controls

import (
	"encoding/json"
	"fmt"
)

// TODO: Make this a generic type with numerical dial as a parameter

type ObservableDialMap struct {
	Observable[ValueMap]
	Dials *NumericDialMap
}

func (m *ObservableDialMap) SetValue(values ValueMap) {
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

// func (m *ObservableDialMap) GetString() string {
// 	r, err := json.Marshal(m)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	return string(r)
// }

func (m *ObservableDialMap) GetValue() ValueMap {
	res := ValueMap{}
	for name, dial := range *m.Dials {
		res[name] = dial.Value
	}
	return res
}

func NewObservableNumericDialMap(channels ...string) *ObservableDialMap {
	return &ObservableDialMap{
		Observable: *NewObservable[ValueMap](),
		Dials:      NewNumericDialMap2(channels...),
	}
}

func (m *ObservableDialMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Dials)
}
