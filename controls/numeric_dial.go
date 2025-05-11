package controls

import (
	"encoding/json"
	"fmt"
	"log"
)

func NewNumericDial() *NumericDial {
	return &NumericDial{
		Value: 0,
	}
}

type NumericDial struct {
	Value byte
}

func (d *NumericDial) SetValue(value byte) {
	d.Value = value
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

func (d *NumericDial) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Value)
}

func (d *NumericDial) SetValueString(value string) {
	var b byte
	if _, err := fmt.Sscanf(value, "%d", &b); err == nil {
		d.SetValue(b)
	} else {
		log.Printf("Error setting value '%s': %s", value, err)
	}
}

func (d *NumericDial) GetValueString() string {
	return fmt.Sprintf("%d", d.Value)
}
