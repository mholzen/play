package controls

import (
	"encoding/json"
	"log"
)

func NewNumericDial() *NumericDial {
	return &NumericDial{
		Value:   0,
		channel: make(chan byte),
	}
}

type NumericDial struct {
	Value   byte
	channel chan byte `json:"-"`
}

func (d *NumericDial) SetValue(value byte) {
	d.Value = value
	d.Emit()
}

func (d *NumericDial) SetMax() {
	d.SetValue(255)
}

func (d *NumericDial) SetMin() {
	d.SetValue(0)
}

func (d *NumericDial) Toggle() {
	if d.Value <= 127 {
		d.SetMax()
	} else {
		d.SetMin()
	}
}

func (d NumericDial) Opposite() byte {
	x := int(d.Value) - 255
	if x < 0 {
		return byte(-x)
	} else {
		return byte(x)
	}
}

func (d *NumericDial) Emit() {
	log.Printf("Emitting %v to %+v", d.Value, d.channel)
	select {
	case d.channel <- d.Value:
	default:
	}
}

func (d *NumericDial) Channel() <-chan byte {
	return d.channel
}

func (d *NumericDial) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Value)
}
