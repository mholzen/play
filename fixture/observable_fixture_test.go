package fixture

import (
	"testing"

	"github.com/mholzen/play-go/controls"
)

func Test_ObservableFixture(t *testing.T) {
	model := NewModelChannels("Foo", []string{"r", "g", "b"})
	var fixture FixtureI = model

	observableFixture := NewObservableFixture(fixture)

	observableFixture.SetChannelValues(controls.ChannelValues{"r": 1})

	var _ FixtureI = observableFixture
}
