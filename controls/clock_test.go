package controls

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_TickPeriod(t *testing.T) {
	c := Clock{
		Bpm:          120.0,
		TicksPerBeat: 4,
		BeatsPerBar:  4,
	}
	assert.Equal(t, 125*time.Millisecond, c.TickPeriod())
}

func TestClock_TicksPerPhrase(t *testing.T) {
	t.Run("calculates ticks per phrase correctly", func(t *testing.T) {
		clock := &Clock{
			TicksPerBeat: 4,
			BeatsPerBar:  4,
			BarPerPhrase: 2,
		}

		expected := 4 * 4 * 2 // 32 ticks per phrase
		actual := clock.TicksPerPhrase()

		assert.Equal(t, expected, actual, "should calculate ticks per phrase correctly")
	})

	t.Run("handles different clock configurations", func(t *testing.T) {
		testCases := []struct {
			name         string
			ticksPerBeat int
			beatsPerBar  int
			barPerPhrase int
			expected     int
		}{
			{"standard 4/4", 4, 4, 4, 64},
			{"waltz 3/4", 4, 3, 4, 48},
			{"compound time", 8, 6, 2, 96},
			{"minimal", 1, 1, 1, 1},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				clock := &Clock{
					TicksPerBeat: tc.ticksPerBeat,
					BeatsPerBar:  tc.beatsPerBar,
					BarPerPhrase: tc.barPerPhrase,
				}

				actual := clock.TicksPerPhrase()
				assert.Equal(t, tc.expected, actual, "should calculate ticks per phrase for %s", tc.name)
			})
		}
	})
}
