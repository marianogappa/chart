package color

import "strings"

var palette = []string{
	`"#f44336"`, `"#e91e63"`, `"#9c27b0"`, `"#673ab7"`, `"#3f51b5"`, `"#2196f3"`, `"#03a9f4"`, `"#00bcd4"`,
	`"#009688"`, `"#4caf50"`, `"#8bc34a"`, `"#cddc39"`, `"#ffeb3b"`, `"#ffc107"`, `"#ff9800"`, `"#ff5722"`,
	`"#795548"`, `"#9e9e9e"`, `"#607d8b"`, `"#f44336"`, `"#e91e63"`, `"#9c27b0"`, `"#673ab7"`, `"#3f51b5"`,
	`"#2196f3"`, `"#03a9f4"`, `"#00bcd4"`, `"#009688"`, `"#4caf50"`, `"#8bc34a"`, `"#cddc39"`, `"#ffeb3b"`,
	`"#ffc107"`, `"#ff9800"`, `"#ff5722"`, `"#795548"`, `"#9e9e9e"`, `"#607d8b"`}

var i = 0

// Reset restarts the palette iterator. Following Reset(), invoking Next() returns the first color in the palette.
func Reset() {
	i = 0
}

// Next iterates through the color palette.
func Next() string {
	result := palette[i]
	i++
	if i > len(palette)-1 {
		i = 0
	}

	return result
}

// FirstN returns a comma-separated string of the first n colors in the palette.
func FirstN(n int) string {
	k := 0
	var cs []string
	for j := 0; j < n; j++ {
		cs = append(cs, palette[k])
		k++
		if k > len(palette)-1 {
			k = 0
		}
	}
	return strings.Join(cs, ",")
}
