package fixture

import "testing"

func Test_ModelChannels(t *testing.T) {
	foo := NewModelChannels("Foo", []string{"r", "g", "b"})

	var _ FixtureI = foo
}
