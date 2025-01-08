package controls

import "testing"

func Test_NumericDialMapIsContainer(t *testing.T) {
	dialMap := NewNumericDialMap("r", "g", "b")
	var _ Container = dialMap
}
