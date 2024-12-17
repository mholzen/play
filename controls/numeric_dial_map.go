package controls

type NumericDialMap map[string]*NumericDial

func NewNumericDialMap(channels ...string) *NumericDialMap {
	res := NumericDialMap{}
	for _, channel := range channels {
		res[channel] = NewNumericDial()
	}
	return &res
}

func (m *NumericDialMap) GetChannels() []string {
	res := []string{}
	for channel := range *m {
		res = append(res, channel)
	}
	return res
}
