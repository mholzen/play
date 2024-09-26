package controls

import (
	"testing"
)

func Test_GetMidiClock(t *testing.T) {
	t.Skip()
	// depends on the midi clock being present and sending

	ticker, err := GetMidiClockTicker()
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 5; i++ {
		select {
		case tick := <-ticker:
			t.Logf("tick: %d", tick)
		default:
		}
	}
}
