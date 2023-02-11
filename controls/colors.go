package controls

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Color struct {
	Red   byte `json:"r"`
	Green byte `json:"g"`
	Blue  byte `json:"b"`
	White byte `json:"w"`
	Amber byte `json:"a"`
	UV    byte `json:"uv"`
}

func (c Color) Values() ValueMap {
	return ValueMap{
		"r":  c.Red,
		"g":  c.Green,
		"b":  c.Blue,
		"w":  c.White,
		"a":  c.Amber,
		"uv": c.UV,
	}
}

type Colors map[string]Color

var AllColors Colors

func init() {
	// Open our jsonFile
	jsonFile, err := os.Open("colors.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &AllColors)
	if err != nil {
		log.Fatal(err)
	}
}
