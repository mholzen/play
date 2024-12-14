package controls

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Toggle struct {
	DiscreteDial[bool] `json:"Value"`
}

func (t *Toggle) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func NewToggle() *Toggle {
	return &Toggle{}
}

func (t *Toggle) On() {
	t.Set(true)
}

func (t *Toggle) Off() {
	t.Set(false)
}

func (t *Toggle) Toggle() {
	if t.Value {
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

func (t *Toggle) SetValueString(value string) {
	value = strings.ToLower(value)
	if value == "true" || value == "1" || value == "on" || value == "yes" || value == "enable" {
		t.Set(true)
	} else {
		t.Set(false)
	}
}

func (t *Toggle) GetValue() bool {
	return t.Value
}

func (t *Toggle) GetValueString() string {
	return fmt.Sprintf("%v", t.Value)
}
