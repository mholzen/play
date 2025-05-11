package controls

type ObservableNumericDial struct {
	NumericDial
	*Observable[byte]
}

func NewObservableNumericDial() *ObservableNumericDial {
	return &ObservableNumericDial{
		NumericDial: NumericDial{
			Value: 0,
		},
		Observable: NewObservable[byte](),
	}
}

func (d *ObservableNumericDial) SetValue(value byte) {
	d.NumericDial.SetValue(value)
	d.Observable.Notify(value)
}

func (d *ObservableNumericDial) SetValueString(value string) {
	d.NumericDial.SetValueString(value)
	d.Observable.Notify(d.Value)
}
