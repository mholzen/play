package patterns

import (
	"testing"

	"github.com/fogleman/ease"
	"github.com/stretchr/testify/assert"
)

func TestEaseSequenceComposition(t *testing.T) {
	sequence := NewEaseSequence(ease.Linear, ease.InQuad, ease.OutQuad)

	result := sequence(0.5)
	assert.IsType(t, float64(0), result)
	assert.GreaterOrEqual(t, result, 0.0)
	assert.LessOrEqual(t, result, 1.0)
}

func TestInvertComposition(t *testing.T) {
	inverted := Invert(ease.Linear)

	result1 := inverted(0.0)
	assert.Equal(t, 1.0, result1)

	result2 := inverted(1.0)
	assert.Equal(t, 0.0, result2)

	result3 := inverted(0.5)
	assert.Equal(t, 0.5, result3)
}

func TestEaseMapComposition(t *testing.T) {
	easeFunc := EaseMap["linear"]
	assert.NotNil(t, easeFunc)

	result := easeFunc(0.5)
	assert.Equal(t, 0.5, result)

	triangleFunc := EaseMap["triangle"]
	assert.NotNil(t, triangleFunc)

	result = triangleFunc(0.5)
	assert.IsType(t, float64(0), result)
}

func TestFunctionChaining(t *testing.T) {
	compose := func(f1, f2 ease.Function) ease.Function {
		return func(t float64) float64 {
			return f2(f1(t))
		}
	}

	chained := compose(ease.Linear, Invert(ease.Linear))
	result := chained(0.3)
	assert.Equal(t, 0.7, result)
}

func TestEaseMapAllFunctions(t *testing.T) {
	testCases := []string{
		"linear", "square", "sine", "sawtooth_up", "sawtooth_down",
		"triangle", "sine_up", "sine_down", "sine_up_down",
		"quad_up", "quad_down", "quad_up_down",
	}

	for _, name := range testCases {
		easeFunc := EaseMap[name]
		assert.NotNil(t, easeFunc, "ease function %s should exist", name)

		result := easeFunc(0.5)
		assert.IsType(t, float64(0), result)
		assert.GreaterOrEqual(t, result, 0.0)
		assert.LessOrEqual(t, result, 1.0)
	}
}

func TestEaseFlat(t *testing.T) {
	result := EaseFlat(0.0)
	assert.Equal(t, 0.0, result)

	result = EaseFlat(0.5)
	assert.Equal(t, 0.0, result)

	result = EaseFlat(1.0)
	assert.Equal(t, 0.0, result)
}

func TestBoundaryValues(t *testing.T) {
	testFunctions := []ease.Function{
		ease.Linear,
		ease.InQuad,
		ease.OutQuad,
		ease.InCubic,
		ease.OutCubic,
	}

	for _, fn := range testFunctions {
		assert.Equal(t, 0.0, fn(0.0))
		assert.Equal(t, 1.0, fn(1.0))
	}
}
