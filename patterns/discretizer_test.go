package patterns

import (
	"fmt"
	"testing"

	"github.com/fogleman/ease"
	"github.com/mholzen/play-go/controls"
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

func ExampleDiscretizer() {
	discretizer := NewDiscretizer(0, 255, ease.InOutQuad, 8)
	values := discretizer.GetValues()

	fmt.Printf("Discretized values: %v\n", values)

	sequence := controls.NewSequence(values)

	fmt.Printf("Sequence values:\n")
	for i := 0; i < len(values); i++ {
		fmt.Printf("Step %d: %d\n", i, sequence.Value())
		sequence.Inc()
	}

	// Output:
	// Discretized values: [0 10 42 94 161 213 245 255]
	// Sequence values:
	// Step 0: 0
	// Step 1: 10
	// Step 2: 42
	// Step 3: 94
	// Step 4: 161
	// Step 5: 213
	// Step 6: 245
	// Step 7: 255
}

func ExampleDiscretizer_square() {
	discretizer := NewDiscretizer(0, 1, ease.Linear, 2)
	values := discretizer.GetValues()
	fmt.Printf("Discretized values: %v\n", values)
	// Output:
	// Discretized values: [0 1]
}

func TestDiscretizerWithSequence(t *testing.T) {
	t.Run("creates sequence compatible slice", func(t *testing.T) {
		discretizer := NewDiscretizer(10, 90, ease.OutBounce, 5)
		values := discretizer.GetValues()

		sequence := controls.NewSequence(values)

		assert.Equal(t, values[0], sequence.Value(), "first sequence value should match first discretized value")

		sequence.Inc()
		assert.Equal(t, values[1], sequence.Value(), "second sequence value should match second discretized value")
	})
}
