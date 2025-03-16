package controls

import "testing"

func TestJoin(t *testing.T) {
	// a join should be an observable
	observable1 := NewObservable[int]()
	observable2 := NewObservable[int]()
	join := NewObservableJoin(observable1, observable2)

	var _ Observable[int] = join

	observable1.Notify(1)
	observable2.Notify(2)
}
