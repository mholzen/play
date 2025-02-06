package controls

import (
	"testing"
)

func Test_DialListHasOrder(t *testing.T) {
	dialMap := NewObservableNumericDialMap("r", "g", "b")
	dialList := NewDialList2(dialMap)

	// Verify dialList implements Container interface
	var _ OrderedContainer = dialList
}
