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

func NewNumericDialMap(channels ...string) *NumericDialMap {
	res := NumericDialMap{}
	for _, channel := range channels {
		res[channel] = NewNumericDial()
	}
	return &res
}
