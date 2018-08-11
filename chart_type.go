package main

import (
	"fmt"

	"github.com/marianogappa/chart/format"
)

func resolveChartType(ct chartType, lf format.LineFormat, datasetLength int) (chartType, error) {
	ct = _resolveChartType(ct, lf)
	return ct, assertChartable(ct, lf, datasetLength)
}

func _resolveChartType(ct chartType, f format.LineFormat) chartType {
	if ct != undefinedChartType {
		return ct
	}
	switch {
	case !f.HasFloats:
		return pie
	case f.FloatCount >= 2 && !f.HasStrings:
		return scatter
	case f.FloatCount > 1:
		return line
	default:
		return pie
	}
}

func assertChartable(ct chartType, f format.LineFormat, datasetLength int) error {
	if datasetLength == 0 {
		return fmt.Errorf("empty dataset; nothing to plot here")
	}
	var errIncompatibleFormat = fmt.Errorf("I don't know how to plot a dataset with this line format")
	switch ct {
	case pie, bar:
		if !f.HasFloats || (f.FloatCount == 1 && !f.HasStrings && !f.HasDateTimes) {
			return errIncompatibleFormat
		}
	case line:
		if !f.HasFloats || (f.FloatCount < 2 && !f.HasStrings && !f.HasDateTimes) {
			return errIncompatibleFormat
		}
	case scatter:
		if !f.HasFloats {
			return errIncompatibleFormat
		}
	}
	return nil
}
