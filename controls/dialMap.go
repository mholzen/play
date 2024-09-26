package controls

type DialMap map[string]*NumericDial

func (m DialMap) SetValues(values ValueMap) {
	for name, value := range values {
		m[name].SetValue(value)
	}
}

func (m DialMap) GetString() string {
	return "<multiple items>"
}
