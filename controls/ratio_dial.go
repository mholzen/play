package controls

import "encoding/json"

type RatioDial = DiscreteDial[Ratio]

type ObservableRatioDial struct {
	RatioDial
	Observers[Ratio]
}

func (d *ObservableRatioDial) MarshalJSON() ([]byte, error) {
	return json.Marshal(&d.RatioDial)
}

func (d *ObservableRatioDial) SetValueString(value string) {
	d.RatioDial.SetValueString(value)
	d.Notify(d.Get())
}

func NewObservableRatioDial() *ObservableRatioDial {
	res := &ObservableRatioDial{
		RatioDial: *NewDiscreteDial(CommonRatios),
		Observers: *NewObservable[Ratio](),
	}
	res.Set(IdentityRatio)
	return res
}

func NewCommonIncreasingRatioDial() *ObservableRatioDial {
	return &ObservableRatioDial{
		RatioDial: *NewDiscreteDial(IncreasingRatios),
		Observers: *NewObservable[Ratio](),
	}
}

func NewCommonDecreasingRatioDial() *ObservableRatioDial {
	return &ObservableRatioDial{
		RatioDial: *NewDiscreteDial(DecreasingRatios),
		Observers: *NewObservable[Ratio](),
	}
}

type ratioJSON struct {
	Value       float64 `json:"value"`
	Numerator   int     `json:"numerator"`
	Denominator int     `json:"denominator"`
	Label       string  `json:"label"`
}

func (r Ratio) MarshalJSON() ([]byte, error) {
	return json.Marshal(ratioJSON{
		Value:       r.ToFloat(),
		Numerator:   r.Numerator,
		Denominator: r.Denominator,
		Label:       r.Label(),
	})
}
