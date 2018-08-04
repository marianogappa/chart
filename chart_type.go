package main

import (
	"github.com/marianogappa/chart/format"
)

func resolveChartType(t chartType, lineFormat format.LineFormat, fss [][]float64, sss [][]string) chartType {
	if t == undefinedChartType {
		switch {
		case !lineFormat.HasFloats:
			return pie
		case lineFormat.FloatCount >= 2 && !lineFormat.HasStrings:
			return scatter
		case lineFormat.FloatCount > 1:
			return line
		default:
			return pie
		}
	}
	return t
}
