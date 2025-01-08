package controls

type Ratio struct {
	numerator   int
	denominator int
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
	if r.denominator == 0 {
		return 0
	}
	return float64(r.numerator) / float64(r.denominator)
}
