package controls

import (
	"testing"
)

func Test_DialListIsContainer(t *testing.T) {
	dialMap := NewObservableNumericDialMap("r", "g", "b")
	dialList := NewDialListFromContainer(dialMap)

	// Verify dialList implements Container interface
	var _ Container = dialList

}
