package controls

import (
	"encoding/json"
	"fmt"
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
	t.SetValue(true)
}

func (t *Toggle) Off() {
	t.SetValue(false)
}

func (t *Toggle) Toggle() {
	if t.Value {
		t.Off()
	} else {
		t.On()
	}
}

func (t *Toggle) GetValue() bool {
	return t.Value
}

func (t *Toggle) GetString() string {
	return fmt.Sprintf("%v", t.Value)
}
