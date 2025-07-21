package controls

import (
	"fmt"
	"strconv"
)

type FloatDial struct {
	Value float64 `json:"value"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
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

func (d *FloatDial) SetValueString(value string) error {
	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	d.Value = parsedValue
	return nil
}

type ObservableFloatDial struct {
	Observers[FloatDial]
	FloatDial
}
