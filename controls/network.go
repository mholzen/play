package controls

import (
	"log"
	"play-go/fixture"

	"github.com/reugn/go-streams"
	ext "github.com/reugn/go-streams/extension"
)

func LinkFixtureChannel(fixture fixture.FixtureI, channelName string, channel <-chan byte) {
	go func() {
		log.Printf("starting listener on %+v", channel)
		for value := range channel {
			log.Printf("received %+v", value)
			fixture.SetValue(channelName, value)
		}
		log.Print("listener ended")
	}()
}

func NewFixtureSink(fixture fixture.FixtureI, channel string) streams.Sink {
	c := make(chan any)

	go func() {
		for v := range c {
			log.Printf("received %+v", v)
			fixture.SetValue(channel, v.(byte))
		}
	}()

	sink := ext.NewChanSink(c)
	return sink
}
