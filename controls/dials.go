package controls

import (
	"fmt"
)

type ChannelList []string

type DialList struct {
	DialMap *ObservableDialMap
	ChannelList
}

func (dl *DialList) SetChannelValue(channel string, value byte) {
	(*dl.DialMap.Dials)[channel].Value = value
}

func (dl *DialList) GetItem(name string) (Item, error) {
	dial, ok := (*dl.DialMap.Dials)[name]
	if !ok {
		return nil, fmt.Errorf("item not found: %s", name)
	}
	return dial, nil
}

func (dl *DialList) Items() map[string]Item {
	items := make(map[string]Item)
	for _, channel := range dl.ChannelList {
		items[channel] = (*dl.DialMap.Dials)[channel]
	}
	return items
}

func (dl *DialList) Keys() []string {
	return dl.ChannelList
}

func (dl *DialList) MarshalJSON() ([]byte, error) {
	return OrderedContainerMarshalJSON(dl)
}

var DefaultChannelList = []string{"r", "g", "b", "a", "w", "uv", "dimmer", "strobe", "speed", "tilt", "pan", "mode"}

func NewDialList(dialMap *ObservableDialMap) *DialList {
	return &DialList{DialMap: dialMap, ChannelList: DefaultChannelList}
}
