package controls

import "fmt"

type NumericDialMap map[string]*NumericDial

func (m *NumericDialMap) GetChannels() []string {
	res := []string{}
	for channel := range *m {
		res = append(res, channel)
	}
	return res
}

func (m *NumericDialMap) GetItem(key string) (Item, error) {
	item, ok := (*m)[key]
	if !ok {
		return nil, fmt.Errorf("item not found")
	}
	return item, nil
}

func (m *NumericDialMap) Items() map[string]Item {
	items := make(map[string]Item)
	for k, v := range *m {
		items[k] = v
	}
	return items
}

func (m *NumericDialMap) SetValue(values ChannelValues) {
	for name, value := range values {
		if dial, ok := (*m)[name]; ok {
			dial.SetValue(value)
		} else {
			panic(fmt.Sprintf("dial '%s' not found", name))
		}
	}
}

func NewNumericDialMap(channels ...string) *NumericDialMap {
	res := NumericDialMap{}
	for _, channel := range channels {
		dial := NewNumericDial()
		res[channel] = dial
	}
	return &res
}
