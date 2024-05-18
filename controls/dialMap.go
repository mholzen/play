package controls

type DialMap map[string]*Dial

func (m DialMap) SetValues(values ValueMap) {
	for name, value := range values {
		m[name].SetValue(value)
	}
}
