package main

import "strings"

func resolveChartType(t chartType, lf string, fss [][]float64, sss [][]string) chartType {
	if t == undefinedChartType {
		switch {
		case strings.Index(lf, "f") == -1:
			return pie
		case strings.Count(lf, "f") > 1:
			return line
		//TODO add more inference cases when there are more chart types available
		default:
			return pie
		}
	}
	return t
}
