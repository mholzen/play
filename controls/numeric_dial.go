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
	Value byte `json:"value"`
	min   byte
	max   byte
}

func (d *NumericDial) Set(value byte) {
	d.Value = value
}

func (d *NumericDial) Get() byte {
	return d.Value
}

func (d *NumericDial) Min() byte {
	return d.min
}

func (d *NumericDial) Max() byte {
	return d.max
}

// Should a NumericDial automatically observable? No, it can have just a polling API
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
	d.SetValue(d.max)
}

func (d *NumericDial) SetMin() {
	d.SetValue(d.min)
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
	NumericDial
}

func (d *ObservableNumericalDial) SetValue(value byte) {
	d.NumericDial.SetValue(value)
	d.Notify(value)
}

func NewObservableNumericalDial(dial *NumericDial) *ObservableNumericalDial {
	return &ObservableNumericalDial{
		Observers:   *NewObservable[byte](),
		NumericDial: *dial,
	}
}

func NewObservableNumericalDial2() Dial[byte] {
	dial := NewNumericDial()
	return NewObservableNumericalDial(dial)
}

func NewObservableNumericalDialFromChannel(channel string) *ObservableNumericalDial {
	return NewObservableNumericalDial(NewNumericDial())
}

func (d *ObservableNumericalDial) GetValueString() string {
	return d.NumericDial.GetValueString()
}

func (d *ObservableNumericalDial) SetValueString(value string) {
	d.NumericDial.SetValueString(value)
	d.Notify(d.NumericDial.Value)
}

func (d *ObservableNumericalDial) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.NumericDial)
}
