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

func (d *NumericDial) SetValue(value byte) {
	d.Value = value
}

func (d *NumericDial) GetValueString() string {
	return fmt.Sprintf("%d", d.Value)
}

func (d *NumericDial) SetValueString(value string) error {
	byteValue, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	d.SetValue(byte(byteValue))
	return nil
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

type ObservableNumericDial struct {
	// NOTE: a NumericDial does not need to automatically be observable
	// as it is useful with simply a polling API
	NumericDial
	Observers[byte]
}

func (d *ObservableNumericDial) Set(value byte) {
	d.NumericDial.SetValue(value)
	d.Notify(value)
}

func (d *ObservableNumericDial) SetValue(value byte) {
	d.NumericDial.SetValue(value)
	d.Notify(value)
}

func (d *ObservableNumericDial) GetValueString() string {
	return d.NumericDial.GetValueString()
}

func (d *ObservableNumericDial) SetValueString(value string) error {
	err := d.NumericDial.SetValueString(value)
	if err != nil {
		return err
	}
	d.Notify(d.Value)
	return nil
}

func (d *ObservableNumericDial) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.NumericDial)
}

func NewObservableNumericDial() *ObservableNumericDial {
	return &ObservableNumericDial{
		Observers:   *NewObservable[byte](),
		NumericDial: *NewNumericDial(),
	}
}

func NewDialObservableNumeric() Dial[byte] {
	// TODO: can this be constructed functionally?
	return NewObservableNumericDial()
}
