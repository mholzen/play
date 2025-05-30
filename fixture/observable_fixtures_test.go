package fixture

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ObservableFixtures2(t *testing.T) {
	model := NewModelChannels("a Model Name", []string{"r", "g", "b"})
	var aFixture Fixture = ChannelFixture{Model: &model, Values: model.GetEmptyValues()}

	fixtures := NewAddressableFixtures[Fixture]()
	fixtures.AddFixture(aFixture, 1)

	observableFixtures := NewIndividualObservableFixtures(fixtures)

	observableFixtures.SetChannelValue("r", 1)
	assert.Equal(t, byte(1), aFixture.GetChannelValue("r"))
}
