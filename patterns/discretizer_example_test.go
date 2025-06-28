package patterns

import (
	"fmt"
	"testing"

	"github.com/fogleman/ease"
	"github.com/mholzen/play-go/controls"
)

func ExampleDiscretizer() {
	discretizer := NewDiscretizer(0, 255, ease.InOutQuad, 8)
	values := discretizer.GetValues()

	fmt.Printf("Discretized values: %v\n", values)

	sequence := controls.NewSequenceT(values)

	fmt.Printf("Sequence values:\n")
	for i := 0; i < len(values); i++ {
		fmt.Printf("Step %d: %d\n", i, sequence.Values())
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

		sequence := controls.NewSequenceT(values)

		if sequence.Values() != values[0] {
			t.Errorf("first sequence value should match first discretized value")
		}

		sequence.Inc()
		if sequence.Values() != values[1] {
			t.Errorf("second sequence value should match second discretized value")
		}
	})
}
