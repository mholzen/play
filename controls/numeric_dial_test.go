package controls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumericDial(t *testing.T) {
	d := NewNumericDial()
	d.SetValue(100)
	assert.Equal(t, byte(100), d.Value)

	var _ Control = d
}
