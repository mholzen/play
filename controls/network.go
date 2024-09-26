package controls

import (
	"log"

	"github.com/mholzen/play-go/fixture"

	"github.com/reugn/go-streams"
	ext "github.com/reugn/go-streams/extension"
)

type LinkFixtureChannel2 struct {
	Fixture     fixture.FixtureI
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
				l.Fixture.SetValue(l.ChannelName, value)
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

func LinkFixtureChannel(fixture fixture.FixtureI, channelName string, channel <-chan byte) *LinkFixtureChannel2 {
	l := &LinkFixtureChannel2{
		Fixture:     fixture,
		ChannelName: channelName,
		Channel:     channel,
	}

	l.Run()

	return l
}

func LinkDialToFixtureChannel(dial *NumericDial, fixture fixture.FixtureI, channel string) *LinkFixtureChannel2 {
	return LinkFixtureChannel(fixture, channel, dial.Channel)
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
