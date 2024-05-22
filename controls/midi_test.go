package controls

import "testing"

func Test_GetMidiClock(t *testing.T) {
	ticker, err := GetMidiClockTicker()
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		select {
		case tick := <-ticker:
			t.Logf("tick: %d", tick)
		}
	}
	t.Fail()
	// if tick != 1 {
	// 	t.Errorf("expected 1, got %d", tick)
	// }
}
