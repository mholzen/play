package fixture

import (
	"fmt"
	"log"

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
