package controls

type ChannelList []string

var DefaultChannelList = ChannelList{"r", "g", "b", "a", "w", "uv", "dimmer", "strobe", "speed", "tilt", "pan", "mode"}

var MotionChannelList = ChannelList{"speed", "tilt", "pan"}

var ColorChannelList = ChannelList{"r", "g", "b", "a", "w", "uv"}

var LightChannelList = ChannelList{"dimmer", "strobe"}

type DialList struct { // TODO: should be "ordered dial map", perhaps simply "dials"
	ItemMap
	// DialMap[byte]
	ChannelList
}

func NewDialList(container Container, channelList ChannelList) *DialList {
	return &DialList{
		ItemMap:     container.Items(),
		ChannelList: channelList,
	}
}

func NewDialListFromContainer(container Container) *DialList {
	return &DialList{ItemMap: container.Items(), ChannelList: DefaultChannelList}
}

func NewDialListFromDialMap(dialMap DialMap[byte]) *DialList {
	return &DialList{ItemMap: dialMap.Items(), ChannelList: DefaultChannelList}
}

func NewDialListFromChannelList(channelList ChannelList) *DialList {
	return &DialList{ItemMap: NewItemMap(), ChannelList: channelList}
}

func (dl *DialList) Add(channel string, dial Dial[byte]) {
	dl.ItemMap[channel] = dial
}

func (dl *DialList) Keys() []string {
	return dl.ChannelList
}

func (dl *DialList) MarshalJSON() ([]byte, error) {
	return OrderedContainerMarshalJSON(dl)
}

func (dl *DialList) SetChannelValue(channel string, value byte) {
	item, ok := dl.ItemMap[channel]
	if !ok {
		panic("channel not found")
	}

	dial, ok := item.(*ObservableNumericalDial)
	if !ok {
		panic("dial is not a *ObservableNumericalDial")
	}

	dial.SetValue(value)
}

func (dl *DialList) GetObservables() map[string]Observable[byte] {
	observables := make(map[string]Observable[byte])
	for channel, item := range dl.ItemMap {
		dial, ok := item.(Observable[byte])
		if !ok {
			continue
		}
		observables[channel] = dial
	}
	return observables
}
