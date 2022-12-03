package main

import (
	"testing"

	"github.com/marianogappa/chart/v4/format"
)

func TestResolveChartType(t *testing.T) {
	tests := []struct {
		name      string
		t         chartType
		lf        string
		expectedT chartType
	}{
		{
			name:      "default case",
			t:         undefinedChartType,
			lf:        "sf",
			expectedT: pie,
		},
		{
			name:      "pie selected; inference ignored",
			t:         pie,
			lf:        "sf",
			expectedT: pie,
		},
		{
			name:      "bar selected; inference ignored",
			t:         bar,
			lf:        "sf",
			expectedT: bar,
		},
		{
			name:      "more than one column of floats, with strings",
			t:         undefinedChartType,
			lf:        "sff",
			expectedT: line,
		},
		{
			name:      "more than one column of floats, without strings",
			t:         undefinedChartType,
			lf:        "ff",
			expectedT: scatter,
		},
	}

	for _, ts := range tests {
		lf, _ := format.NewLineFormat(ts.lf, ' ', "") // ignoring errors as we're not testing the format package here
		result, err := resolveChartType(ts.t, lf, 123)
		if err != nil { // TODO test cases where there's an error
			t.Errorf("%v: there was an error resolving the chart type", ts.name)
		}

		if result != ts.expectedT {
			t.Errorf("%v: %v was not equal to %v", ts.name, result, ts.expectedT)
		}
	}
}
