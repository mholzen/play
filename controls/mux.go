package controls

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"slices"
)

type Mux[T any] struct {
	Sources map[string]Emitter[T]
	Source  string
}

func NewMux[T any]() Mux[T] {
	res := Mux[T]{}
	res.Sources = make(map[string]Emitter[T], 0)
	return res
}

func (m *Mux[T]) Add(name string, source Emitter[T]) {
	m.Sources[name] = source
	if len(m.Sources) == 1 {
		m.Source = name
	}
}

func (m *Mux[T]) GetSource() string {
	return m.Source
}

func (m *Mux[T]) SetSource(name string) error {
	if _, ok := m.Sources[name]; !ok {
		return fmt.Errorf("cannot find source '%s'", name)
	}
	m.Source = name
	return nil
}

func (m Mux[T]) GetValue() T {
	return m.Sources[m.Source].GetValue()
}

func (m Mux[T]) Channel() <-chan T {
	ch := make(chan T)
	go func() {
		cases := make([]reflect.SelectCase, 0, len(m.Sources))
		sourceNames := make([]string, 0, len(m.Sources))

		for name, source := range m.Sources {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(source.Channel()),
			})
			sourceNames = append(sourceNames, name)
		}

		for {
			log.Printf("listening on %s", m.Source)

			chosen, value, ok := reflect.Select(cases)
			log.Printf("chosen: %d, value: %v, ok: %t", chosen, value, ok)
			if !ok {
				continue
			}
			sourceName := sourceNames[chosen]
			if sourceName == m.Source {
				log.Printf("emitting %v to %+v", value.Interface().(T), ch)
				ch <- value.Interface().(T)
			} else {
				log.Printf("ignoring %v from %s", value.Interface().(T), sourceName)
			}
		}
	}()
	return ch
}

func (m Mux[T]) MarshalJSON() ([]byte, error) {
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

func (m Mux[T]) GetValueString() string {
	return m.Source
}

func (m *Mux[T]) SetValueString(value string) {
	m.SetSource(value)
}
