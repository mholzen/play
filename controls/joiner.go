package controls

type Joiner[T any] struct {
	Sources []Observable[T]
	Observers[T]
}

func NewJoiner[T any]() *Joiner[T] {
	res := &Joiner[T]{}
	res.Sources = make([]Observable[T], 0)
	res.Observers = *NewObservable[T]()
	return res
}

func (j *Joiner[T]) Add(source Observable[T]) {
	j.Sources = append(j.Sources, source)
	channel := make(chan T)
	source.AddObserver(channel)
	go func() {
		for value := range channel {
			j.Notify(value)
		}
	}()
}
