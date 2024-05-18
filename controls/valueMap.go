package controls

import (
	"fmt"
	"image/color"
	"math/rand"
	"sort"
	"strings"

	"github.com/mholzen/play-go/fixture"
)

type ValueMap map[string]byte

func (values ValueMap) ApplyTo(f fixture.FixtureI) {
	// log.Printf("Applying %s", values.String())
	for k, v := range values {
		f.SetValue(k, v)
	}
}

func (values ValueMap) String() string {
	res := []string{}
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		res = append(res, fmt.Sprintf("%s:%03d", key, values[key]))
	}
	return strings.Join(res, ", ")
}

func InterpolateValues(v1, v2 ValueMap, t float64) ValueMap {
	res := make(ValueMap)
	for k, v := range v1 {
		res[k] = byte(float64(v) + (float64(v2[k])-float64(v))*t)
	}
	return res
}

func NewMapFromColor(c color.RGBA) ValueMap {
	return ValueMap{
		"r":  c.R,
		"g":  c.G,
		"b":  c.B,
		"w":  0,
		"a":  c.A,
		"uv": 0,
	}
}

type ParamValue struct {
	Param string
	Value byte
}

func parseRandomValue(value string) (byte, error) {
	var num1, num2 int
	_, err := fmt.Sscanf(value, "%d-%d", &num1, &num2)
	if err != nil {
		return 0, err
	}
	return byte(rand.Intn(num2-num1) + num1), nil
}

func parseParam(param string) (ParamValue, error) {
	parts := strings.Split(param, ":")
	if len(parts) != 2 {
		return ParamValue{}, fmt.Errorf("param '%s' has %d parts", param, len(parts))
	}
	if strings.ContainsRune(parts[1], '-') {
		v, err := parseRandomValue(parts[1])
		if err != nil {
			return ParamValue{}, err
		}
		return ParamValue{Param: parts[0], Value: v}, nil
	}
	return ParamValue{Param: parts[0], Value: byte(parts[1][0])}, nil
}

func NewMap(param ...string) (ValueMap, error) {
	res := make(ValueMap)
	for _, p := range param {
		v, err := parseParam(p)
		if err != nil {
			return nil, err
		}
		res[v.Param] = v.Value
	}
	return res, nil
}
