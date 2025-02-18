package controls

import (
	"encoding/json"
	"fmt"
)

type Ratio struct {
	Numerator   int `json:"numerator"`
	Denominator int `json:"denominator"`
}

var WellKnownRatios = []Ratio{
	{1, 1},
	{3, 4},
	{2, 3},
	{1, 2},
	{1, 3},
	{1, 4},
	{1, 8},
	{1, 12},
	{1, 16},
	{1, 32},
}

func (r Ratio) ToFloat() float64 {
	if r.Denominator == 0 {
		return 0
	}
	return float64(r.Numerator) / float64(r.Denominator)
}

func (r Ratio) Label() string {
	return fmt.Sprintf("%d:%d", r.Numerator, r.Denominator)
}

type ObservableRatioDial struct {
	DiscreteDial[Ratio]
	Observers[Ratio]
}

func NewObservableRatioDial() *ObservableRatioDial {
	return &ObservableRatioDial{
		DiscreteDial: *NewDiscreteDial(WellKnownRatios),
		Observers:    *NewObservable[Ratio](),
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
