package main

import (
	"flag"
	_ "flag"
	"fmt"
	"log"
	"math"
)
import tsize "github.com/kopoli/go-terminal-size"

type vec2 struct{ x, y float64 }

// Palette source:http://mewbies.com/geek_fun_files/ascii/ascii_art_light_scale_and_gray_scale_chart.htm
type Palette []byte

var MinimalPalette = Palette{
	' ', '.', ',', ':', ';', 'i', 't', '%', 'X', '$', '@', '#',
}

var DefaultPalette = Palette{
	'$', '@', 'B', '%', '8', '&', 'W', 'M', '#', '*', 'o', 'a', 'h', 'k', 'b', 'd', 'p', 'q', 'w',
	'm', 'Z', 'O', '0', 'Q', 'L', 'C', 'J', 'U', 'Y', 'X', 'z', 'c', 'v', 'u', 'n', 'x', 'r', 'j',
	'f', 't', '/', '|', '(', ')', '1', '{', '}', '[', ']', '?', '-', '_', '+', '~', '<', '>', 'i',
	'!', 'l', 'I', ';', ':', ',', '\\', '^', '`', '\'', '.', ' ',
}

var iterations int = 32
var zoom float64 = 0.5
var position vec2 = vec2{x: 0.5, y: 0}
var invert bool = false

func (p Palette) translate(value uint32) byte {
	normalizedValue := float64(value) / float64(iterations)

	if invert {
		normalizedValue = 1.0 - normalizedValue
	}
	maxValue := len(p) - 1
	index := int(math.Round(float64(maxValue) * normalizedValue))
	return p[index]
}

func calculate(p vec2) uint32 {
	z := vec2{x: 0, y: 0}

	for i := 0; i < iterations; i++ {
		z.x = p.x + (z.x*z.x - z.y*z.y)
		z.y = p.y + (2.0 * z.x * z.y)

		dot := z.x*z.x + z.y*z.y

		// Number does not belong in the set
		if dot > 4.0 {
			return uint32(i)
		}
	}

	return 0
}

func main() {
	flag.Float64Var(&zoom, "zoom", 1.0, "zoom factor")
	flag.Float64Var(&position.x, "x", 0.0, "position x")
	flag.Float64Var(&position.y, "y", 0.0, "position y")
	flag.IntVar(&iterations, "iterations", 128, "number of iterations")
	flag.BoolVar(&invert, "invert", false, "invert palette")
	flag.Parse()

	var s tsize.Size

	s, err := tsize.GetSize()
	if err != nil {
		log.Fatal("Failed to retrieve terminal dimensions")
	}

	buffer := make([]uint32, s.Width*s.Height)

	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			index := y*s.Width + x
			w := float64(s.Width)
			h := float64(s.Height)
			aspectRatio := w / h
			magnification := 1.0 / zoom
			normX := ((((float64(x)/w)*2.0 - 1.0) * aspectRatio) + position.x) * magnification
			normY := (((float64(y)/h)*2.0 - 1.0) + position.y) * magnification
			buffer[index] = calculate(vec2{x: normX, y: normY})
		}
	}

	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			index := y*s.Width + x
			val := buffer[index]
			c := DefaultPalette.translate(val)
			fmt.Print(string(c))
		}
		fmt.Printf("\n")
	}
}
