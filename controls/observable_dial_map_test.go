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
