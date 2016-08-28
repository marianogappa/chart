package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestColorNextCycles(t *testing.T) {
	colorReset()
	color := colorNext()
	for i := 0; i < len(colorPalette)-1; i++ {
		colorNext()
	}
	newColor := colorNext()

	if color != newColor {
		t.Error("the colorNext function doesn't seem to cycle after finishing available colors")
	}
}

func TestColorReset(t *testing.T) {
	colorReset()
	color := colorNext()
	for i := 0; i < len(colorPalette)/2; i++ {
		colorNext()
	}
	colorReset()
	newColor := colorNext()

	if color != newColor {
		t.Errorf("the colorReset() function doesn't seem to restart the color palette iterator; mismatch: %v != %v", color, newColor)
	}
}

func TestColorFirstN(t *testing.T) {
	for i := 0; i < len(colorPalette)-1; i++ {
		cs := strings.Split(colorFirstN(i), ",")
		if i > 0 && len(cs) != i {
			t.Errorf("the colorFirstN() function returned %v colors when asked for %v colors", len(cs), i)
		}
		d := map[string]struct{}{}
		for _, c := range cs {
			d[c] = struct{}{}
		}
		if len(d) < len(cs) {
			t.Errorf("the colorFirstN(%v) function contained %v duplicates! %v != %v", i, len(cs)-len(d), cs, d)
		}
		if i >= 2 {
			if !reflect.DeepEqual(cs[:i-1], strings.Split(colorFirstN(i-1), ",")) {
				t.Errorf("the colorFirstN(%v) function returned different initial colours than colorFirstN(%v-1)", i, i)
			}
		}
	}
}
