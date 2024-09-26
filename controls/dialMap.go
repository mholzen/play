package controls

import "encoding/json"

type DialMap map[string]*NumericDial

func (m DialMap) SetValues(values ValueMap) {
	for name, value := range values {
		m[name].SetValue(value)
	}
}

func (m DialMap) GetString() string {
	r, err := json.Marshal(m)
	if err != nil {
		return err.Error()
	}
	return string(r)
}
