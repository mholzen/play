package patterns

import (
	"testing"

	"github.com/fogleman/ease"
)

func TestDiscretizer(t *testing.T) {
	t.Run("linear interpolation", func(t *testing.T) {
		d := NewDiscretizer(0, 100, ease.Linear, 5)
		values := d.GetValues()

		expected := []int{0, 25, 50, 75, 100}
		if len(values) != len(expected) {
			t.Fatalf("expected %d values, got %d", len(expected), len(values))
		}

		for i, v := range values {
			if v != expected[i] {
				t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
			}
		}
	})

	t.Run("single interval", func(t *testing.T) {
		d := NewDiscretizer(10, 20, ease.Linear, 1)
		values := d.GetValues()

		if len(values) != 1 {
			t.Fatalf("expected 1 value, got %d", len(values))
		}
		if values[0] != 10 {
			t.Errorf("expected 10, got %d", values[0])
		}
	})

	t.Run("zero intervals", func(t *testing.T) {
		d := NewDiscretizer(0, 100, ease.Linear, 0)
		values := d.GetValues()

		if len(values) != 0 {
			t.Fatalf("expected 0 values, got %d", len(values))
		}
	})

	t.Run("with easing function", func(t *testing.T) {
		d := NewDiscretizer(0, 100, ease.InQuad, 3)
		values := d.GetValues()

		if len(values) != 3 {
			t.Fatalf("expected 3 values, got %d", len(values))
		}

		if values[0] != 0 {
			t.Errorf("first value should be 0, got %d", values[0])
		}
		if values[2] != 100 {
			t.Errorf("last value should be 100, got %d", values[2])
		}

		if values[1] >= 50 {
			t.Errorf("middle value with InQuad should be less than linear interpolation (50), got %d", values[1])
		}
	})
}
