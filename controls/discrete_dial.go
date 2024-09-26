package controls

type DiscreteDial[T any] struct {
	Value T
	C     chan T `json:"-"`
}

func (d *DiscreteDial[T]) SetValue(value T) {
	d.Value = value
	d.Emit()
}

func (d *DiscreteDial[T]) Emit() {
	select {
	case d.C <- d.Value:
	default:
	}
}

func (d *DiscreteDial[T]) Channel() <-chan T {
	return d.C
}
