package fixture

import "testing"

func Test_ObservableFixtures2(t *testing.T) {
	fixtures := NewFixturesGeneric[Fixture]()
	observableFixtures := NewObservableFixtures(fixtures)

	var _ AddressableFixtures[Fixture] = observableFixtures.AddressableFixtures
}
