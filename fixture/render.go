package fixture

import (
	"log"
	"time"

	"github.com/akualab/dmx"
)

const REFRESH = 40 * time.Millisecond // DMXIS cannot read faster than 40ms

func Render(fixtures Fixtures[Fixture], connection dmx.DMX) error {
	ticker := time.NewTicker(REFRESH)
	go func() {
		for range ticker.C {
			values := fixtures.GetByteArray()
			// log.Printf("values: %v", values)
			for address, value := range values {
				if address == 0 {
					continue
				}
				connection.SetChannel(address, value)
			}
			err := connection.Render()
			if err != nil {
				log.Printf("ERROR rendering error: %s", err)
			}
		}
	}()
	return nil
}
