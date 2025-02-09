package controls

type ChannelList []string

var DefaultChannelList = []string{"r", "g", "b", "a", "w", "uv", "dimmer", "strobe", "speed", "tilt", "pan", "mode"}

type DialList struct {
	ItemMap
	ChannelList
}

func NewDialList(container Container) *DialList {
	return &DialList{ItemMap: container.Items(), ChannelList: DefaultChannelList}
}

func (dl *DialList) Keys() []string {
	return dl.ChannelList
}

func (dl *DialList) MarshalJSON() ([]byte, error) {
	return OrderedContainerMarshalJSON(dl)
}
