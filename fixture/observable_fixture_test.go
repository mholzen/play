package fixture

import (
	"testing"

	"github.com/mholzen/play-go/controls"
)

func Test_ObservableFixture(t *testing.T) {
	foo := NewModelChannels("Foo", []string{"r", "g", "b"})
	var fixture FixtureI = foo

	observableFixture := NewObservableFixture(fixture)

	observableFixture.SetValueMap(controls.ValueMap{"r": 1})

	var _ FixtureI = observableFixture
}