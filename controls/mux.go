package controls

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
)

type Mux[T any] struct {
	Sources map[string]Observable[T]
	Source  string
	Observers[T]
}

func NewMux[T any]() *Mux[T] {
	res := &Mux[T]{}
	res.Sources = make(map[string]Observable[T], 0)
	res.Observers = *NewObservable[T]()
	return res
}

func (m *Mux[T]) Add(name string, source Observable[T]) {
	m.Sources[name] = source
	if len(m.Sources) == 1 {
		m.Source = name
	}
	channel := make(chan T)
	source.AddObserver(channel)
	go func() {
		for value := range channel {
			if name == m.Source {
				// log.Printf("mux value received from %s: notifying", name)
				m.Notify(value)
			} else {
				// log.Printf("mux value received from %s: ignoring", name)
			}
		}
	}()
}

func (m *Mux[T]) GetSource() string {
	return m.Source
}

func (m *Mux[T]) SetSource(name string) error {
	if _, ok := m.Sources[name]; !ok {
		return fmt.Errorf("cannot find source '%s'", name)
	}
	log.Printf("mux setting source to %s", name)
	m.Source = name
	return nil
}

func (m *Mux[T]) MarshalJSON() ([]byte, error) {
	res := struct {
		Options []string `json:"options"`
		Value   string   `json:"value"`
	}{
		Options: make([]string, 0, len(m.Sources)),
		Value:   m.Source,
	}

	for name := range m.Sources {
		res.Options = append(res.Options, name)
	}
	slices.Sort(res.Options)

	return json.Marshal(res)
}

func (m *Mux[T]) GetValueString() string {
	return m.Source
}

func (m *Mux[T]) SetValueString(value string) {
	m.SetSource(value)
}
