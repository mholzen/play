package controls

import "testing"

func Test_DialMap_IsContainer(t *testing.T) {
	var _ Container = &DialMap[byte]{}
}
