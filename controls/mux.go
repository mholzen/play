package controls

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
)

type Mux[T any] struct {
	Sources map[string]ObservableI[T]
	Source  string
	Observable[T]
}

func NewMux[T any]() *Mux[T] {
	res := &Mux[T]{}
	res.Sources = make(map[string]ObservableI[T], 0)
	res.Observable = *NewObservable[T]()
	return res
}

func (m *Mux[T]) Add(name string, source ObservableI[T]) {
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
		Sources []string `json:"sources"`
		Source  string   `json:"source"`
	}{
		Sources: make([]string, 0, len(m.Sources)),
		Source:  m.Source,
	}

	for name := range m.Sources {
		res.Sources = append(res.Sources, name)
	}
	slices.Sort(res.Sources)

	return json.Marshal(res)
}

func (m *Mux[T]) GetValueString() string {
	return m.Source
}

func (m *Mux[T]) SetValueString(value string) {
	m.SetSource(value)
}