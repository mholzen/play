package controls

import (
	"log"
	"testing"
	"time"

	ext "github.com/reugn/go-streams/extension"
	"github.com/reugn/go-streams/flow"
)

func Test_flow(t *testing.T) {

	source := ext.NewChanSource(tickerChan(time.Millisecond * 500))
	flow := flow.NewPassThrough()
	c := make(chan any)
	sink := ext.NewChanSink(c)

	go func() {
		log.Print("listening")
		for v := range c {
			log.Printf("received %+v", v)
		}
	}()

	source.Via(flow).To(sink)

	time.Sleep(2 * time.Second)
}

func tickerChan(repeat time.Duration) chan any {
	ticker := time.NewTicker(repeat)
	oc := ticker.C
	nc := make(chan any)
	go func() {
		for t := range oc {
			log.Printf("sending %+v", t)
			nc <- 1
		}
	}()
	return nc
}
