package main

import "strings"

var colorPalette = []string{
	`"#f44336"`, `"#e91e63"`, `"#9c27b0"`, `"#673ab7"`, `"#3f51b5"`, `"#2196f3"`, `"#03a9f4"`, `"#00bcd4"`,
	`"#009688"`, `"#4caf50"`, `"#8bc34a"`, `"#cddc39"`, `"#ffeb3b"`, `"#ffc107"`, `"#ff9800"`, `"#ff5722"`,
	`"#795548"`, `"#9e9e9e"`, `"#607d8b"`}

var colorI = 0

// Reset restarts the palette iterator. Following Reset(), invoking Next() returns the first color in the palette.
func colorReset() {
	colorI = 0
}

// Next iterates through the color palette.
func colorNext() string {
	result := colorPalette[colorI]
	colorI++
	if colorI > len(colorPalette)-1 {
		colorI = 0
	}

	return result
}

// FirstN returns a comma-separated string of the first n colors in the palette.
func colorFirstN(n int) string {
	k := 0
	var cs []string
	for j := 0; j < n; j++ {
		cs = append(cs, colorPalette[k])
		k++
		if k > len(colorPalette)-1 {
			k = 0
		}
	}
	return strings.Join(cs, ",")
}
