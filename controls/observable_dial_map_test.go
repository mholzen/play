package controls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObservableDialMap(t *testing.T) {
	s := NewObservableNumericDialMap("r", "g", "b")

	var _ Container = s
}

func Test_ObservableDialMapIsContainer(t *testing.T) {
	dialMap := NewObservableNumericDialMap("r", "g", "b")

	// Verify dialMap implements Container interface
	var _ Container = dialMap

	json, err := dialMap.MarshalJSON()
	assert.Nil(t, err)
	assert.Contains(t, string(json), `"r":0`)
	assert.Contains(t, string(json), `"g":0`)
	assert.Contains(t, string(json), `"b":0`)
}

func Test_ObservableDialMap2(t *testing.T) {
	dialMap := NewObservableDialMap2()
	dialMap.AddItem("r", NewObservableNumericDial())

	item, err := dialMap.GetItem("r")
	assert.Nil(t, err)

	dial, ok := item.(*ObservableNumericDial)
	assert.True(t, ok)

	ch := make(chan ChannelValues)
	dialMap.AddObserver(ch)

	advance := make(chan int)

	go func() {
		dial.SetValue(100)
		<-advance // wait for 1
	}()

	values := <-ch
	assert.Equal(t, values["r"], byte(100))
	advance <- 1
}
