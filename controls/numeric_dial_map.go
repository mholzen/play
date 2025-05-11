package controls

import "fmt"

type ObservableNumericalDialmap map[string]*ObservableNumericDial

func NewObservableNumericDialMap2(channels ...string) *ObservableNumericalDialmap {
	res := ObservableNumericalDialmap{}
	for _, channel := range channels {
		dial := NewObservableNumericDial()
		res[channel] = dial
	}
	return &res
}

func (m *ObservableNumericalDialmap) GetItem(name string) (Item, error) {
	dial, ok := (*m)[name]
	if !ok {
		return nil, fmt.Errorf("dial '%s' not found", name)
	}
	return dial, nil
}

func (m *ObservableNumericalDialmap) Items() map[string]Item {
	res := make(map[string]Item)
	for name, dial := range *m {
		res[name] = dial
	}
	return res
}

func (m *ObservableNumericalDialmap) GetChannelValues() ChannelValues {
	res := ChannelValues{}
	for name, dial := range *m {
		res[name] = dial.Value
	}
	return res
}
