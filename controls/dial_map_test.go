package controls

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DialMap(t *testing.T) {
	dialMap := NewNumericDialMap("ch1", "ch2")

	dialMap.SetValue(ValueMap{"ch1": 42, "ch2": 24})

	require.Equal(t, ValueMap{"ch1": 42, "ch2": 24}, dialMap.GetValue())
}
