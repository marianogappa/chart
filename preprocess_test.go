package main

import (
	"reflect"
	"testing"
)

func TestPreprocess(t *testing.T) {
	tests := []struct {
		name      string
		i         []string
		o         options
		fss       [][]float64
		sss       [][]string
		expectedO options
	}{
		{
			name:      "empty case",
			i:         []string{},
			o:         options{separator: '\t', scaleType: linear, chartType: pie},
			fss:       [][]float64{},
			sss:       [][]string{},
			expectedO: options{separator: '\t', scaleType: linear, chartType: pie},
		},
		{
			name: "just strings",
			i:    []string{"a", "b", "c", "a"},
			o:    options{separator: '\t', scaleType: linear, chartType: pie},
			fss: [][]float64{
				[]float64{2}, []float64{1}, []float64{1},
			},
			sss: [][]string{
				[]string{"a"}, []string{"b"}, []string{"c"},
			},
			expectedO: options{separator: '\t', scaleType: linear, chartType: pie},
		},
	}

	for _, ts := range tests {
		fss, sss, o := preprocess(ts.i, ts.o)
		if !reflect.DeepEqual(fss, ts.fss) {
			t.Errorf("'%v' failed: %v was not equal to %v", ts.name, fss, ts.fss)
		}
		if !reflect.DeepEqual(sss, ts.sss) {
			t.Errorf("'%v' failed: %v was not equal to %v", ts.name, sss, ts.sss)
		}
		if !reflect.DeepEqual(o, ts.o) {
			t.Errorf("'%v' failed: %v was not equal to %v", ts.name, o, ts.expectedO)
		}
	}
}
