package controls

import (
	"log"
	"time"

	"github.com/nsf/termbox-go"
)

type TickEventInterface interface {
	BPM() float64
}

type TickEvent = time.Time

type TickerChannelInterface interface {
	GetChannel() <-chan time.Time
}

type TerminalTickerChannel struct {
	ch chan TickEvent
}

func NewTerminalTickerChannel() *TerminalTickerChannel {
	tc := &TerminalTickerChannel{
		ch: make(chan TickEvent),
	}
	go tc.start()
	return tc
}

func (tc *TerminalTickerChannel) start() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlC {
				close(tc.ch)
				return
			}
			log.Printf("=== TAP")
			tc.ch <- time.Now()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func (tc *TerminalTickerChannel) GetChannel() <-chan TickEvent {
	return tc.ch
}
