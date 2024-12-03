package controls

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DialMap(t *testing.T) {
	dialMap := NewNumericDialMap("test")

	received := make(chan ValueMap)
	go func() {
		value := <-dialMap.Channel()
		received <- value
	}()

	dialMap.SetValue(ValueMap{"test": 42})

	value := <-received
	require.Equal(t, ValueMap{"test": 42}, value)
}
