package controls

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MultiCounter(t *testing.T) {
	mc := MultiCounter{
		Counter{Period: 4},
		Counter{Period: 3},
		Counter{Period: 2},
	}

	for i := 0; i < 10; i++ {
		mc.Inc()
	}
	require.Equal(t, []int{2, 2, 0}, mc.Values())

	for i := 0; i < 10; i++ {
		mc.Inc()
	}
	require.Equal(t, []int{0, 2, 1}, mc.Values())

	for i := 0; i < 4; i++ {
		mc.Inc()
	}
	require.Equal(t, []int{0, 0, 0}, mc.Values())
}
