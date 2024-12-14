package controls

type NumericDialMap map[string]*NumericDial

func NewNumericDialMap2(channels ...string) *NumericDialMap {
	res := NumericDialMap{}
	for _, channel := range channels {
		res[channel] = NewNumericDial()
	}
	return &res
}
