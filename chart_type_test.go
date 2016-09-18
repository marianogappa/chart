package main

import "testing"

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
			lf:        "s,f",
			fss:       [][]float64{},
			sss:       [][]string{},
			expectedT: pie,
		},
		{
			name:      "pie selected; inference ignored",
			t:         pie,
			lf:        "s,f",
			fss:       [][]float64{},
			sss:       [][]string{},
			expectedT: pie,
		},
		{
			name:      "bar selected; inference ignored",
			t:         bar,
			lf:        "s,f",
			fss:       [][]float64{},
			sss:       [][]string{},
			expectedT: bar,
		},
		{
			name:      "more than one column of floats, with strings",
			t:         undefinedChartType,
			lf:        "s,f,f",
			fss:       [][]float64{},
			sss:       [][]string{},
			expectedT: line,
		},
		{
			name:      "more than one column of floats, without strings",
			t:         undefinedChartType,
			lf:        "f,f",
			fss:       [][]float64{},
			sss:       [][]string{},
			expectedT: scatter,
		},
	}

	for _, ts := range tests {
		result := resolveChartType(ts.t, ts.lf, ts.fss, ts.sss)

		if result != ts.expectedT {
			t.Errorf("%v: %v was not equal to %v", ts.name, result, ts.expectedT)
		}
	}
}
