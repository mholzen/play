package controls

import "testing"

func Test_RatioDialIsControl(t *testing.T) {
	var _ Control = &RatioDial{}
}

func Test_ObservableRatioDialIsControl(t *testing.T) {
	var _ Control = &ObservableRatioDial{}
}
