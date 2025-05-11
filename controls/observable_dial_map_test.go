package controls

import (
	"testing"
)

func TestObservableDialMap(t *testing.T) {
	s := NewObservableNumericDialMap()

	var _ Container = s
}
