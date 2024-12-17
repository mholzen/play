package fixture

import (
	"github.com/mholzen/play-go/controls"
)

func LinkObservableToFixture(source controls.Observable[FixtureValues], target *Fixtures[Fixture]) {
	channel := make(chan FixtureValues)
	source.AddObserver(channel)
	go func() {
		for fixtureValues := range channel {
			// log.Printf("mux output: %v", fixtureValues)
			(*target).SetValue(fixtureValues)
		}
	}()
}
