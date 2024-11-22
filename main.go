package main

import (
	"flag"
	"fmt"
	"log"
	"math"
)
import tsize "github.com/kopoli/go-terminal-size"

type vec2 struct{ x, y float64 }

// Palette source:http://mewbies.com/geek_fun_files/ascii/ascii_art_light_scale_and_gray_scale_chart.htm
type Palette []byte

var DefaultPalette = Palette{
	'$', '@', 'B', '%', '8', '&', 'W', 'M', '#', '*', 'o', 'a', 'h', 'k', 'b', 'd', 'p', 'q', 'w',
	'm', 'Z', 'O', '0', 'Q', 'L', 'C', 'J', 'U', 'Y', 'X', 'z', 'c', 'v', 'u', 'n', 'x', 'r', 'j',
	'f', 't', '/', '|', '(', ')', '1', '{', '}', '[', ']', '?', '-', '_', '+', '~', '<', '>', 'i',
	'!', 'l', 'I', ';', ':', ',', '\\', '^', '`', '\'', '.', ' ',
}

var iterations = 32
var zoom = 0.5
var position = vec2{x: 0.5, y: 0}
var invert = false
var aspectCorrection = true

func (p Palette) translate(value uint32, currMin uint32, currMax uint32) byte {

	valueRange := currMax - currMin

	normalizedValue := float64(value-currMin) / float64(valueRange)

	if invert {
		normalizedValue = 1.0 - normalizedValue
	}

	paletteLength := len(p) - 1
	index := int(math.Round(float64(paletteLength) * normalizedValue))
	return p[index]
}

func findMaxMin(arr []uint32) (uint32, uint32) {
	// Initialize the variables to hold the maximum and minimum values to draw comparisons.
	currMax := arr[0]
	currMin := arr[0]
	// Iterate over the array
	for i := 1; i < len(arr); i++ {
		// if the current element is greater than the present maximum
		if arr[i] > currMax {
			currMax = arr[i]
		}
		// if the current element is smaller than the present minimum
		if arr[i] < currMin {
			currMin = arr[i]
		}
	}

	if currMin == currMax {
		currMax = currMin + 1
	}

	return currMin, currMax
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
	flag.BoolVar(&aspectCorrection, "correct_aspect", false, "apply aspect correction")
	flag.Parse()

	var s tsize.Size

	s, err := tsize.GetSize()
	if err != nil || s.Width == 0 || s.Height == 0 {
		log.Fatal("Failed to retrieve terminal dimensions")
	}

	buffer := make([]uint32, s.Width*s.Height)

	var aspectRatio = 1.0
	w := float64(s.Width)
	h := float64(s.Height)
	magnification := 1.0 / zoom

	if aspectCorrection {
		/* The constant factor below accounts for average terminal character dimensions ratio
		** TODO: auto detect this factor, current method is likely to produce inaccurate results depending on font.
		 */
		aspectRatio = 0.65 * w / h
	}

	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			index := y*s.Width + x
			normX := ((((float64(x)/w)*2.0 - 1.0) * aspectRatio) + position.x) * magnification
			normY := (((float64(y)/h)*2.0 - 1.0) + position.y) * magnification
			buffer[index] = calculate(vec2{x: normX, y: normY})
		}
	}

	currMin, currMax := findMaxMin(buffer)

	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			index := y*s.Width + x
			val := buffer[index]
			c := DefaultPalette.translate(val, currMin, currMax)
			fmt.Print(string(c))
		}
		fmt.Printf("\n")
	}
}
