package controls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DialListIsContainer(t *testing.T) {
	dialMap := NewObservableNumericDialMap("r", "g", "b")
	dialList := NewDialList(dialMap)

	// Verify dialList implements Container interface
	var _ Container = dialList

	json, err := dialList.MarshalJSON()
	assert.Nil(t, err)
	assert.Contains(t, string(json), `{"name":"r","value":0}`)
}
