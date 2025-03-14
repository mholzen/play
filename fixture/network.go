package fixture

import (
	"log"

	"github.com/mholzen/play-go/controls"
)

func ConnectObservableChannelValuesToFixtures(source controls.Observable[controls.ChannelValues], target Fixtures[Fixture]) {
	ch := make(chan controls.ChannelValues)
	source.AddObserver(ch)
	go func() {
		for value := range ch {
			SetChannelValues(target, value)
		}
	}()
}

func ConnectObservableToFixtures(source controls.Observable[byte], channel string, target Fixtures[Fixture]) {
	ch := make(chan byte)
	source.AddObserver(ch)
	go func() {
		for value := range ch {
			target.SetChannelValue(channel, value)
			log.Printf("set %s to %d in Connect fixture", channel, value)
		}
	}()
}

func ConnectObservablesToFixtures(sources map[string]controls.Observable[byte], target Fixtures[Fixture]) {
	for channel, dial := range sources {
		ConnectObservableToFixtures(dial, channel, target)
	}
}

func ConnectObservableValuesToFixtures(source controls.Observable[FixtureValues], target Fixtures[Fixture]) {
	channel := make(chan FixtureValues)
	source.AddObserver(channel)
	go func() {
		for fixtureValues := range channel {
			// log.Printf("mux output: %v", fixtureValues)
			target.SetValue(fixtureValues)
		}
	}()
}

func ConnectDialMapToFixtures(dialMap controls.Observable[controls.ChannelValues], target Fixtures[Fixture]) {
	received := make(chan controls.ChannelValues)
	dialMap.AddObserver(received)
	go func() {
		for channelValues := range received {
			for channel, value := range channelValues {
				target.SetChannelValue(channel, value)
			}
		}
	}()
}
