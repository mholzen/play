package controls

import (
	"fmt"
	"testing"
)

type TestLabelable int

func (t TestLabelable) Label() string {
	return fmt.Sprintf("%d", t)
}

func Test_DiscreteDial_Is_Control(t *testing.T) {
	var _ Control = &DiscreteDial[TestLabelable]{}
}
