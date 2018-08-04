package main

import (
	"testing"

	"github.com/marianogappa/chart/format"
)

func TestResolveChartType(t *testing.T) {
	tests := []struct {
		name      string
		t         chartType
		lf        string
		fss       [][]float64
		sss       [][]string
		expectedT chartType
	}{
		{
			name:      "default case",
			t:         undefinedChartType,
			lf:        "sf",
			fss:       [][]float64{},
			sss:       [][]string{},
			expectedT: pie,
		},
		{
			name:      "pie selected; inference ignored",
			t:         pie,
			lf:        "sf",
			fss:       [][]float64{},
			sss:       [][]string{},
			expectedT: pie,
		},
		{
			name:      "bar selected; inference ignored",
			t:         bar,
			lf:        "sf",
			fss:       [][]float64{},
			sss:       [][]string{},
			expectedT: bar,
		},
		{
			name:      "more than one column of floats, with strings",
			t:         undefinedChartType,
			lf:        "sff",
			fss:       [][]float64{},
			sss:       [][]string{},
			expectedT: line,
		},
		{
			name:      "more than one column of floats, without strings",
			t:         undefinedChartType,
			lf:        "ff",
			fss:       [][]float64{},
			sss:       [][]string{},
			expectedT: scatter,
		},
	}

	for _, ts := range tests {
		lf, _ := format.NewLineFormat(ts.lf, ' ', "") // ignoring errors as we're not testing the format package here
		result := resolveChartType(ts.t, lf, ts.fss, ts.sss)

		if result != ts.expectedT {
			t.Errorf("%v: %v was not equal to %v", ts.name, result, ts.expectedT)
		}
	}
}
