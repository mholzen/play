package fixture

import (
	"testing"

	"github.com/mholzen/play/controls"
	"github.com/stretchr/testify/assert"
)

func Test_ObservableFixture(t *testing.T) {
	model := NewModelChannelsWithName("a Model Name", []string{"r", "g", "b"})
	var fixture Fixture = ChannelFixture{Model: &model, Values: model.GetEmptyValues()}

	observableFixture := NewObservableFixture(fixture)

	observableFixture.SetChannelValues(controls.ChannelValues{"r": 1})

	assert.Equal(t, byte(1), fixture.GetChannelValue("r"))

	var _ Fixture = observableFixture
}
