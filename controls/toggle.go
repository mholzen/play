package controls

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ToggleLabelable bool

func (t ToggleLabelable) Label() string {
	return fmt.Sprintf("%v", t)
}

type Toggle struct {
	DiscreteDial[ToggleLabelable] `json:"Value"`
}

func (t *Toggle) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Get())
}

func NewToggle() *Toggle {
	return &Toggle{
		DiscreteDial: *NewDiscreteDial([]ToggleLabelable{false, true}),
	}
}

func (t *Toggle) On() {
	t.Set(true)
}

func (t *Toggle) Off() {
	t.Set(false)
}

func (t *Toggle) Toggle() {
	if t.Get() {
		t.Off()
	} else {
		t.On()
	}
}

func (t *Toggle) SetValue(value bool) {
	if value {
		t.Set(true)
	} else {
		t.Set(false)
	}
}

func (t *Toggle) SetValueString(value string) error {
	value = strings.ToLower(value)
	if value == "true" || value == "1" || value == "on" || value == "yes" || value == "enable" {
		t.Set(true)
	} else {
		t.Set(false)
	}
	return nil
}

func (t *Toggle) GetValue() bool {
	return bool(t.DiscreteDial.Get())
}

func (t *Toggle) GetValueString() string {
	return fmt.Sprintf("%v", t.DiscreteDial.Get())
}

type ObservableToggle struct {
	Observers[Toggle]
	Toggle
}
