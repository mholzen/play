package controls

import (
	"encoding/json"
	"fmt"
	"slices"
)

type Selector[T any] struct {
	Options  map[string]T
	Selected string
	Observers[T]
}

func NewSelector[T any]() *Selector[T] {
	res := &Selector[T]{}
	res.Options = make(map[string]T, 0)
	res.Observers = *NewObservable[T]()
	return res
}

func (s *Selector[T]) SetOptions(options map[string]T) {
	s.Options = options
}

func (s *Selector[T]) GetSelected() string {
	return s.Selected
}

func (s *Selector[T]) GetSelectedValue() T {
	return s.Options[s.Selected]
}

func (s *Selector[T]) SetSelected(name string) error {
	if _, ok := s.Options[name]; !ok {
		err := fmt.Errorf("cannot find source '%s'", name)
		panic(err)
	}
	s.Selected = name
	return nil
}

func (s *Selector[T]) MarshalJSON() ([]byte, error) {
	res := struct {
		Options []string `json:"options"`
		Value   string   `json:"value"`
	}{
		Options: make([]string, 0, len(s.Options)),
		Value:   s.Selected,
	}

	for name := range s.Options {
		res.Options = append(res.Options, name)
	}
	slices.Sort(res.Options)

	return json.Marshal(res)
}

func (s *Selector[T]) GetValueString() string {
	return s.Selected
}

func (s *Selector[T]) SetValueString(value string) error {
	err := s.SetSelected(value)
	if err != nil {
		return err
	}
	s.Notify(s.GetSelectedValue())
	return nil
}
