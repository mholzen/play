package controls

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Observable(t *testing.T) {
	wait := make(chan int)
	var expected int

	observable := NewObservable[int]()
	OnChange(observable, func(value int) {
		fmt.Println("value received:", value)
		expected = value
		wait <- 1
	})

	observable.Notify(1)

	<-wait
	require.Equal(t, expected, 1)
}
