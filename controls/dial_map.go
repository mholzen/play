package controls

type DialMap map[string]*NumericDial

func NewDialMap() *DialMap {
	return &DialMap{}
}

func (d *DialMap) SetChannelValue(channel string, value byte) {
	(*d)[channel].SetValue(value)
}
