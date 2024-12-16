package controls

import (
	"sync"
)

type Observable[T any] interface {
	AddObserver(observer chan T)
	RemoveObserver(observer chan T)
	Notify(event T)
}

type Observers[T any] struct {
	observers map[chan T]struct{}
	lock      sync.Mutex
}

func NewObservable[T any]() *Observers[T] {
	return &Observers[T]{
		observers: make(map[chan T]struct{}),
	}
}

func (o *Observers[T]) AddObserver(observer chan T) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.observers[observer] = struct{}{}
}

func (o *Observers[T]) RemoveObserver(observer chan T) {
	o.lock.Lock()
	defer o.lock.Unlock()
	delete(o.observers, observer)
	close(observer) // Close the channel to notify the observer.
}

func (o *Observers[T]) Notify(event T) {
	o.lock.Lock()
	defer o.lock.Unlock()
	for observer := range o.observers {
		observer <- event
	}
}

func (o *Observers[T]) AddObserverFunc(observer func(T)) {
	ch := make(chan T)
	go func() {
		for event := range ch {
			observer(event)
		}
	}()
}
