package patterns

import (
	"testing"

	"github.com/fogleman/ease"
	"github.com/stretchr/testify/assert"
)

func TestDiscretizer(t *testing.T) {
	t.Run("linear interpolation", func(t *testing.T) {
		d := NewDiscretizer(0, 100, ease.Linear, 5)
		values := d.GetValues()

		expected := []int{0, 25, 50, 75, 100}
		assert.Equal(t, expected, values, "should produce correct linear interpolation values")
	})

	t.Run("single interval", func(t *testing.T) {
		d := NewDiscretizer(10, 20, ease.Linear, 1)
		values := d.GetValues()

		assert.Len(t, values, 1, "should return exactly one value")
		assert.Equal(t, 10, values[0], "single interval should return start value")
	})

	t.Run("zero intervals", func(t *testing.T) {
		d := NewDiscretizer(0, 100, ease.Linear, 0)
		values := d.GetValues()

		assert.Empty(t, values, "zero intervals should return empty slice")
	})

	t.Run("with easing function", func(t *testing.T) {
		d := NewDiscretizer(0, 100, ease.InQuad, 3)
		values := d.GetValues()

		assert.Len(t, values, 3, "should return correct number of values")
		assert.Equal(t, 0, values[0], "first value should be start value")
		assert.Equal(t, 100, values[2], "last value should be end value")
		assert.Less(t, values[1], 50, "middle value with InQuad should be less than linear interpolation")
	})
}
