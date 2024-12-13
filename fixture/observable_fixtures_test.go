package fixture

import "testing"

func Test_ObservableFixtures2(t *testing.T) {
	fixtures := NewFixturesGeneric[FixtureI]()
	observableFixtures := NewObservableFixtures(fixtures)

	var _ FixturesGeneric[FixtureI] = observableFixtures.FixturesGeneric
}