package fixture

import (
	"testing"

	"github.com/mholzen/play-go/controls"
)

func Test_ObservableFixture(t *testing.T) {
	model := NewModelChannels("Foo", []string{"r", "g", "b"})
	var fixture Fixture = model

	observableFixture := NewObservableFixture(fixture)

	observableFixture.SetChannelValues(controls.ChannelValues{"r": 1})

	var _ Fixture = observableFixture
}
