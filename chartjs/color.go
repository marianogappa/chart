package chartjs

import "strings"

var colorPalette = [][]string{
	{
		`"#ffbb00"`, `"#f90085"`, `"#00bfb2"`, `"#008ffc"`, `"#fc4f00"`, `"#9d00ff"`, `"#ff0000"`, `"#00b515"`,
		`"#c7f400"`, `"#f9a2a2"`, `"#007072"`, `"#e6beff"`, `"#aa5b00"`, `"#fff266"`, `"#7f0000"`, `"#aaffc3"`,
		`"#7f7f00"`, `"#ffe0c1"`, `"#000080"`, `"#808080"`, `"#000000"`,
		`"#ffe08b"`, `"#fc8bc7"`, `"#8be1dc"`, `"#8bccfd"`, `"#fdaf8b"`, `"#dba2ff"`, `"#ff7373"`, `"#a2e4a9"`,
		`"#eafba2"`, `"#fcdddd"`, `"#a2cbcb"`, `"#f3e1ff"`, `"#e0c3a2"`, `"#fffac7"`, `"#d0a2a2"`, `"#d8ffe3"`,
		`"#d0d0a2"`, `"#fff3e8"`, `"#a2a2d0"`, `"#d0d0d0"`,
		`"#a37700"`, `"#9f0055"`, `"#007a72"`, `"#005ca1"`, `"#a13300"`, `"#56008c"`, `"#8c0000"`, `"#00630c"`,
		`"#5b6f00"`, `"#885959"`, `"#003e3f"`, `"#7e688c"`, `"#5d3200"`, `"#746e2f"`, `"#460000"`, `"#4e7459"`,
		`"#3a3a00"`, `"#746658"`, `"#000052"`},
	{
		`"#f44336"`, `"#9c27b0"`, `"#3f51b5"`, `"#03a9f4"`, `"#009688"`, `"#8bc34a"`, `"#ffeb3b"`, `"#ff9800"`,
		`"#795548"`, `"#607d8b"`, `"#e91e63"`, `"#673ab7"`, `"#2196f3"`, `"#00bcd4"`, `"#4caf50"`, `"#cddc39"`,
		`"#ffc107"`, `"#ff5722"`, `"#9e9e9e"`},
	{
		`"#08306b"`, `"#08519c"`, `"#2171b5"`, `"#4292c6"`, `"#6baed6"`, `"#9ecae1"`, `"#c6dbef"`, `"#deebf7"`,
		`"#00441b"`, `"#006d2c"`, `"#238b45"`, `"#41ab5d"`, `"#74c476"`, `"#a1d99b"`, `"#c7e9c0"`, `"#e5f5e0"`,
		`"#7f2704"`, `"#a63603"`, `"#d94801"`, `"#f16913"`, `"#fd8d3c"`, `"#fdae6b"`, `"#fdd0a2"`, `"#fee6ce"`,
		`"#3f007d"`, `"#54278f"`, `"#6a51a3"`, `"#807dba"`, `"#9e9ac8"`, `"#bcbddc"`, `"#dadaeb"`, `"#efedf5"`,
		`"#67001f"`, `"#980043"`, `"#ce1256"`, `"#e7298a"`, `"#df65b0"`, `"#c994c7"`, `"#d4b9da"`, `"#e7e1ef"`,
		`"#000000"`, `"#252525"`, `"#525252"`, `"#737373"`, `"#969696"`, `"#bdbdbd"`, `"#d9d9d9"`, `"#f0f0f0"`},
}

var colorI = 0

// Reset restarts the palette iterator. Following Reset(), invoking Next() returns the first color in the palette.
func colorReset() {
	colorI = 0
}

// Next iterates through the color palette.
func colorNext(i int) string {
	result := colorPalette[i][colorI]
	colorI++
	if colorI > len(colorPalette[i])-1 {
		colorI = 0
	}

	return result
}

func colorIndex(i, j int) string {
	return colorPalette[i][j%len(colorPalette[i])]
}

// FirstN returns a comma-separated string of the first n colors in the palette.
func colorFirstN(i, n int) string {
	k := 0
	var cs []string
	for j := 0; j < n; j++ {
		cs = append(cs, colorPalette[i][k])
		k++
		if k > len(colorPalette[i])-1 {
			k = 0
		}
	}
	return strings.Join(cs, ",")
}

func colorRepeat(i, j, n int) string {
	return strings.Repeat(colorIndex(i, j)+",", n)
}
