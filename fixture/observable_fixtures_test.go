package fixture

import "testing"

func Test_ObservableFixtures2(t *testing.T) {
	fixtures := NewAddressableFixtures[Fixture]()
	observableFixtures := NewObservableFixtures(fixtures)

	var _ AddressableFixtures[Fixture] = observableFixtures.AddressableFixtures
}
