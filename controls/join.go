package controls

type ObservableJoin[T any] struct {
	observables []Observable[T]
}

func NewObservableJoin[T any](observables ...Observable[T]) *ObservableJoin[T] {
	return &ObservableJoin[T]{observables: observables}
}

func (o *ObservableJoin[T]) AddObserver(observer chan T) {
	for _, observable := range o.observables {
		observable.AddObserver(observer)
	}
}

func (o *ObservableJoin[T]) RemoveObserver(observer chan T) {
	for _, observable := range o.observables {
		observable.RemoveObserver(observer)
	}
}

func (o *ObservableJoin[T]) Notify(event T) {
	for _, observable := range o.observables {
		observable.Notify(event)
	}
}
