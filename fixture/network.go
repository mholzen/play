package fixture

import (
	"log"

	"github.com/mholzen/play-go/controls"
)

func LinkObservableToFixture(source controls.ObservableI[FixtureValues], target *FixturesInterface[FixtureI]) {
	channel := make(chan FixtureValues)
	source.AddObserver(channel)
	go func() {
		for fixtureValues := range channel {
			log.Printf("mux output: %v", fixtureValues)
			(*target).SetValue(fixtureValues)
		}
	}()
}
