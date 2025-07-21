package controls

type ObservableNumericDial struct { // TODO: reconcile with ObservableNumericalDial
	NumericDial
	Observable[byte]
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

func (d *ObservableNumericDial) SetValueString(value string) error {
	err := d.NumericDial.SetValueString(value)
	if err != nil {
		return err
	}
	d.Notify(d.Value)
	return nil
}
