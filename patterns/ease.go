package patterns

import (
	"github.com/fogleman/ease"
)

func NewEaseSequence(easeFuncs ...ease.Function) ease.Function {
	return func(t float64) float64 {
		index := int(t * float64(len(easeFuncs)))
		if index >= len(easeFuncs) {
			index = len(easeFuncs) - 1
		}
		return easeFuncs[index](t)
	}
}

func Invert(easeFunc ease.Function) ease.Function {
	return func(t float64) float64 {
		return 1 - easeFunc(t)
	}
}

func EaseFlat(t float64) float64 {
	return 0
}

var EaseMap = map[string]ease.Function{
	"linear":          ease.Linear,
	"square":          ease.InOutSquare,
	"sine":            ease.InOutSine,
	"sawtooth_up":     NewEaseSequence(ease.Linear),
	"sawtooth_down":   NewEaseSequence(Invert(ease.Linear)),
	"triangle":        NewEaseSequence(ease.Linear, Invert(ease.Linear)),
	"sine_up":         ease.InSine,
	"sine_down":       ease.OutSine,
	"sine_up_down":    ease.InOutSine,
	"quad_up":         ease.InQuad,
	"quad_down":       ease.OutQuad,
	"quad_up_down":    ease.InOutQuad,
	"cubic_up":        ease.InCubic,
	"cubic_down":      ease.OutCubic,
	"cubic_up_down":   ease.InOutCubic,
	"quart_up":        ease.InQuart,
	"quart_down":      ease.OutQuart,
	"quart_up_down":   ease.InOutQuart,
	"quint_up":        ease.InQuint,
	"quint_down":      ease.OutQuint,
	"quint_up_down":   ease.InOutQuint,
	"expo_up":         ease.InExpo,
	"expo_down":       ease.OutExpo,
	"expo_up_down":    ease.InOutExpo,
	"circ_up":         ease.InCirc,
	"circ_down":       ease.OutCirc,
	"circ_up_down":    ease.InOutCirc,
	"back_up":         ease.InBack,
	"back_down":       ease.OutBack,
	"back_up_down":    ease.InOutBack,
	"elastic_up":      ease.InElastic,
	"elastic_down":    ease.OutElastic,
	"elastic_up_down": ease.InOutElastic,
	"bounce_up":       ease.InBounce,
	"bounce_down":     ease.OutBounce,
	"bounce_up_down":  ease.InOutBounce,
}
