package fixture

import (
	"fmt"
	"log"
	"time"

	"github.com/akualab/dmx"
)

var devices = []string{"/dev/ttyUSB0", "/dev/ttyUSB1", "/dev/tty.usbserial-ENVVVCOF"}

func GetConnection() (*dmx.DMX, error) {
	for _, device := range devices {
		connection, err := dmx.NewDMXConnection(device)
		if err != nil {
			log.Printf("Could not connect to DMX device at %s: %s -- skipping", device, err)
			continue
		}
		if connection != nil {
			log.Printf("Connected to DMX device at %s", device)
			return connection, nil
		}
	}
	return nil, fmt.Errorf("could not connect to any DMX device")
}

const REFRESH = 40 * time.Millisecond // DMXIS cannot read faster than 40ms

func Render(f Fixtures2, connection *dmx.DMX) {
	ticker := time.NewTicker(REFRESH)
	go func() {
		for range ticker.C {
			f.Render(*connection)
		}
	}()
}
