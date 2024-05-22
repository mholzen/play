package main

import (
	"log"
	"testing"

	"go.bug.st/serial"
)

func Test_Open(t *testing.T) {
	dev := "/dev/tty.usbserial-ENVVVCOF"
	// dev := "/dev/cu.usbserial-ENVVVCOF"
	mode := &serial.Mode{
		BaudRate: 250000,
		DataBits: 8,
		StopBits: serial.TwoStopBits,
		Parity:   serial.NoParity,
	}

	// From enttec dmx usb pro manual:
	// 'baudRate': 250000,
	// 'dataBits': 8,
	// 'stopBits': 2,
	// 'parity': 'none',

	log.Printf("opening serial")
	serial, err := serial.Open(dev, mode)
	if err != nil {
		log.Printf("error opening serial port: %s", err)
		return
	}
	log.Printf("open serial: %v", serial)
}
