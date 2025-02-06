package controls

import (
	"fmt"
	"testing"
)

func Test_Observable(t *testing.T) {
	observable := NewObservable[int]()
	observable.AddObserverFunc(func(value int) {
		fmt.Println("value changed", value)
	})

	observable.Notify(1)
	observable.Notify(2)
	observable.Notify(3)
	observable.Notify(4)
	observable.Notify(5)
}
