package controls

import (
	"sync"
)

type ObservableI[T any] interface {
	AddObserver(observer chan T)
	RemoveObserver(observer chan T)
	Notify(event T)
}

type Observable[T any] struct {
	observers map[chan T]struct{}
	lock      sync.Mutex
}

func NewObservable[T any]() *Observable[T] {
	return &Observable[T]{
		observers: make(map[chan T]struct{}),
	}
}

func (o *Observable[T]) AddObserver(observer chan T) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.observers[observer] = struct{}{}
}

func (o *Observable[T]) RemoveObserver(observer chan T) {
	o.lock.Lock()
	defer o.lock.Unlock()
	delete(o.observers, observer)
	close(observer) // Close the channel to notify the observer.
}

func (o *Observable[T]) Notify(event T) {
	o.lock.Lock()
	defer o.lock.Unlock()
	for observer := range o.observers {
		// log.Printf("notifying observer")
		observer <- event
		// select {
		// case
		// // Event sent successfully.
		// observer <- event:
		// default:
		// 	// Drop the event if the Observer channel is full.
		// }
	}
}

func (o *Observable[T]) AddObserverFunc(observer func(T)) {
	ch := make(chan T)
	go func() {
		for event := range ch {
			observer(event)
		}
	}()
}
