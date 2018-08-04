package chartjs

import (
	"reflect"
	"strings"
	"testing"
)

func TestColorNextCycles(t *testing.T) {
	for i := 0; i < len(colorPalette); i++ {
		colorReset()
		color := colorNext(i)
		for j := 0; j < len(colorPalette[i])-1; j++ {
			colorNext(i)
		}
		newColor := colorNext(i)

		if color != newColor {
			t.Error("the colorNext function doesn't seem to cycle after finishing available colors")
		}
	}
}

func TestColorReset(t *testing.T) {
	for i := 0; i < len(colorPalette); i++ {
		colorReset()
		color := colorNext(i)
		for j := 0; j < len(colorPalette[i])/2; j++ {
			colorNext(i)
		}
		colorReset()
		newColor := colorNext(i)

		if color != newColor {
			t.Errorf("the colorReset() function doesn't seem to restart the color palette iterator; mismatch: %v != %v", color, newColor)
		}
	}
}

func TestColorFirstN(t *testing.T) {
	for i := 0; i < len(colorPalette); i++ {
		for j := 0; j < len(colorPalette[i])-1; j++ {
			cs := strings.Split(colorFirstN(i, j), ",")
			if j > 0 && len(cs) != j {
				t.Errorf("the colorFirstN() function returned %v colors when asked for %v colors", len(cs), j)
			}
			d := map[string]struct{}{}
			for _, c := range cs {
				d[c] = struct{}{}
			}
			if len(d) < len(cs) {
				t.Errorf("the colorFirstN(%v) function contained %v duplicates! %v != %v", j, len(cs)-len(d), cs, d)
			}
			if j >= 2 {
				if !reflect.DeepEqual(cs[:j-1], strings.Split(colorFirstN(i, j-1), ",")) {
					t.Errorf("the colorFirstN(%v) function returned different initial colours than colorFirstN(%v-1)", j, j)
				}
			}
		}
	}
}
