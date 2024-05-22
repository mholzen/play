//go:build darwin

package controls

import (
	"log"
	"strings"

	"github.com/youpy/go-coremidi"
)

type TickerChannel <-chan int

var CLOCK = byte(0xF8)
var SYNC = byte(0xFA)

func GetMidiClockTicker() (TickerChannel, error) {
	client, err := coremidi.NewClient("midi clock client")
	if err != nil {
		return nil, err
	}

	tick := 0
	ch := make(chan int)

	port, err := coremidi.NewInputPort(
		client,
		"test",
		func(source coremidi.Source, packet coremidi.Packet) {
			if !strings.Contains(source.Name(), "Traktor") {
				return
			}
			log.Printf("packet: %v", packet)

			if packet.Data[0] == CLOCK {
				tick++
				ch <- tick
			}
			if packet.Data[0] == SYNC {
				tick = 0
				ch <- tick
			}
		},
	)
	if err != nil {
		return nil, err
	}

	sources, err := coremidi.AllSources()
	if err != nil {
		return nil, err
	}

	for _, source := range sources {
		if !strings.Contains(source.Name(), "Traktor") {
			continue
		}
		// func(source coremidi.Source) {
		port.Connect(source)
		// }(source)
	}

	return ch, nil
}

// func main() {
// 	client, err := coremidi.NewClient("a client")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	var CLOCK = byte(0xF8)
// 	var SYNC = byte(0xFA)

// 	var tick int
// 	port, err := coremidi.NewInputPort(
// 		client,
// 		"test",
// 		func(source coremidi.Source, packet coremidi.Packet) {
// 			if !strings.Contains(source.Name(), "Traktor") {
// 				return
// 			}

// 			if packet.Data[0] == CLOCK {
// 				tick++
// 			}
// 			if packet.Data[0] == SYNC {
// 				tick = 0
// 			}
// 			if tick%24 == 0 {
// 				log.Printf("beat: %d\n", tick)
// 			}
// 		},
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	sources, err := coremidi.AllSources()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	for _, source := range sources {
// 		if !strings.Contains(source.Name(), "Traktor") {
// 			continue
// 		}
// 		log.Printf("source %s", source.Name())
// 		func(source coremidi.Source) {
// 			port.Connect(source)
// 		}(source)
// 	}

// 	ch := make(chan int)
// 	<-ch
// }
