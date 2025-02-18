package controls

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func NewNumericDial() *NumericDial {
	return &NumericDial{
		Value: 0,
	}
}

type Settable interface {
	SetValue(value byte)
}

type NumericDial struct {
	Value byte
}

// Should a NumericDial automatically observable? No.  It can have just a polling interface.

func (d *NumericDial) SetValue(value byte) {
	d.Value = value
}

func (d *NumericDial) GetValueString() string {
	return fmt.Sprintf("%d", d.Value)
}

func (d *NumericDial) SetValueString(value string) {
	byteValue, err := strconv.Atoi(value)
	if err != nil {
		return
	}
	d.SetValue(byte(byteValue))
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

type ObservableNumericalDial struct {
	Observers[byte]
	Dial *NumericDial // TODO: can be embedded?
}

func (d *ObservableNumericalDial) SetValue(value byte) {
	d.Dial.SetValue(value)
	d.Notify(value)
}

func NewObservableNumericalDial(dial *NumericDial) *ObservableNumericalDial {
	return &ObservableNumericalDial{
		Observers: *NewObservable[byte](),
		Dial:      dial,
	}
}

func (d *ObservableNumericalDial) GetValueString() string {
	return d.Dial.GetValueString()
}

func (d *ObservableNumericalDial) SetValueString(value string) {
	d.Dial.SetValueString(value)
	d.Notify(d.Dial.Value)
}

func (d *ObservableNumericalDial) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Dial)
}
