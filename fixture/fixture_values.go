package fixture

import "github.com/mholzen/play-go/controls"

type FixtureValues map[int]controls.ValueMap

func NewFixtureValues(fixtures []*InstalledFixture) FixtureValues {
	values := make(FixtureValues)
	for _, fixture := range fixtures {
		values[fixture.Address] = fixture.GetValueMap()
	}
	return values
}
