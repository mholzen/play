package controls

import (
	"fmt"
)

type Ratio struct {
	Numerator   int `json:"numerator"`
	Denominator int `json:"denominator"`
}

func (r Ratio) Inverse() Ratio {
	return Ratio{r.Denominator, r.Numerator}
}

type RatioSlice []Ratio

func (s RatioSlice) Reverse() {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func (s RatioSlice) Invert() {
	for i := range s {
		s[i] = s[i].Inverse()
	}
}

var IdentityRatio = Ratio{1, 1}

var DecreasingRatios = RatioSlice{
	{3, 4},
	{2, 3},
	{1, 2},
	{1, 3},
	{1, 4},
	{1, 8},
	{1, 12},
	{1, 16},
	{1, 24},
	{1, 32},
}

var IncreasingRatios, CommonRatios RatioSlice

// var CommonRatios = make([]Ratio, len(IncreasingRatios)+1+len(DecreasingRatios))

func init() {
	IncreasingRatios = make(RatioSlice, len(DecreasingRatios))
	copy(IncreasingRatios, DecreasingRatios)
	IncreasingRatios.Invert()

	CommonRatios = make(RatioSlice, len(DecreasingRatios))
	copy(CommonRatios, DecreasingRatios)
	CommonRatios.Reverse()
	CommonRatios = append(CommonRatios, IdentityRatio)
	CommonRatios = append(CommonRatios, IncreasingRatios...)
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
