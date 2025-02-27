package controls

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DialListIsContainer(t *testing.T) {
	dialMap := NewObservableNumericDialMap("r", "g", "b")
	dialList := NewDialListFromContainer(dialMap)

	// Verify dialList implements Container interface
	var _ Container = dialList

	json, err := dialList.MarshalJSON()
	assert.Nil(t, err)
	assert.Contains(t, string(json), `["r",0]`)
}

func Test_NewDialListFromChannels(t *testing.T) {
	channels := []string{"r", "g", "b"}
	dial := NewObservableNumericalDial2()
	var _ Dial[byte] = dial
	dialMap := ChannelsToDialMap2(channels, NewObservableNumericalDial2)

	bytes, err := json.Marshal(dialMap)
	assert.Nil(t, err)
	assert.Contains(t, string(bytes), `"r":{"value":0}`)

	dialList := NewDialListFromDialMap(dialMap)
	bytes, err = json.Marshal(dialList)
	assert.Nil(t, err)
	assert.Contains(t, string(bytes), `["r",{"value":0}]`)
}
