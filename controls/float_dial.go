package controls

import (
	"fmt"
	"strconv"
)

type FloatDial struct {
	Value float64
	Min   float64
	Max   float64
}

func (d *FloatDial) SetValue(value float64) {
	d.Value = value
}

func (d *FloatDial) SetMin() {
	d.SetValue(d.Min)
}

func (d *FloatDial) SetMax() {
	d.SetValue(d.Max)
}

func (d *FloatDial) GetValueString() string {
	return fmt.Sprintf("%v", d.Value)
}

func (d *FloatDial) SetValueString(value string) {
	d.Value, _ = strconv.ParseFloat(value, 64)
}

type ObservableFloatDial struct {
	Observers[FloatDial]
	FloatDial
}
