package controls

import "encoding/json"

type ChannelList []string

type DialList struct {
	DialMap *ObservableDialMap
	ChannelList
}

// marshalJSON is a custom JSON marshaller for DialList
func (dl DialList) MarshalJSON() ([]byte, error) {
	type DialListItem struct {
		Channel string
		Value   *byte
	}
	res := make([]DialListItem, 0)
	for _, channel := range dl.ChannelList {
		// could account for spaces in channel names here
		item := DialListItem{channel, nil}
		dial, ok := (*dl.DialMap.Dials)[channel]
		if ok {
			item.Value = &dial.Value
		}

		res = append(res, item)
	}
	return json.Marshal(res)
}
