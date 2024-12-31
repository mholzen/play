package controls

type DiscreteDial[T any] struct {
	Value T
}

func (d *DiscreteDial[T]) Set(value T) {
	d.Value = value
}

