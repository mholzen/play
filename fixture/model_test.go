package fixture

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ModelChannels(t *testing.T) {
	model := NewModelChannelsWithName("Foo", []string{"r", "g", "b"})

	assert.Equal(t, []string{"r", "g", "b"}, model.GetChannels())

	// var _ Fixture = model
}
