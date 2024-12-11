package fixture

import (
	"log"

	"github.com/mholzen/play-go/controls"
	"github.com/reugn/go-streams"
	ext "github.com/reugn/go-streams/extension"
)

type LinkFixtureChannel2 struct {
	Fixture     FixtureI
	ChannelName string
	Channel     <-chan byte
	Paused      bool
}

func (l LinkFixtureChannel2) Run() {
	go func() {
		log.Printf("starting listener on %+v", l.Channel)
		for value := range l.Channel {
			log.Printf("received %+v", value)
			if !l.Paused {
				l.Fixture.SetChannelValue(l.ChannelName, value)
			}
		}
		log.Print("listener ended")
	}()
}

func (l *LinkFixtureChannel2) Pause() {
	l.Paused = true
}

func (l *LinkFixtureChannel2) Resume() {
	l.Paused = false
}

func LinkFixtureChannel(fixture FixtureI, channelName string, channel <-chan byte) *LinkFixtureChannel2 {
	l := &LinkFixtureChannel2{
		Fixture:     fixture,
		ChannelName: channelName,
		Channel:     channel,
	}

	l.Run()

	return l
}

func LinkDialToFixtureChannel(dial *controls.NumericDial, fixture FixtureI, channel string) *LinkFixtureChannel2 {
	return LinkFixtureChannel(fixture, channel, dial.Channel())
}

func NewFixtureSink(fixture FixtureI, channel string) streams.Sink {
	c := make(chan any)

	go func() {
		for v := range c {
			log.Printf("received %+v", v)
			fixture.SetChannelValue(channel, v.(byte))
		}
	}()

	sink := ext.NewChanSink(c)
	return sink
}

func LinkObservableToFixture(source controls.ObservableI[FixtureValues], target FixturesInterface[FixtureI]) {
	channel := make(chan FixtureValues)
	source.AddObserver(channel)
	go func() {
		for fixtureValues := range channel {
			log.Printf("mux output received")
			log.Printf("fixtureValues: %d fixtures", len(fixtureValues))
			// log.Printf(": %d fixtures", len(fixtureValues))
			// for address, values := range fixtureValues {
			// 	log.Printf("address: %d, values: %v", address, values)
			// }
			target.SetValue(fixtureValues)
		}
	}()
}
