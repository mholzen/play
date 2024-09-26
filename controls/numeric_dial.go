package controls

import "log"

func NewNumericDial() *NumericDial {
	return &NumericDial{
		Value:   0,
		Channel: make(chan byte),
	}
}

type NumericDial struct {
	Value   byte
	Channel chan byte `json:"-"`
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
	log.Printf("Emitting to %+v", d.Channel)
	select {
	case d.Channel <- d.Value:
	default:
	}
}
