package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestPreprocess(t *testing.T) {
	tests := []struct {
		name      string
		i         []string
		o         options
		fss       [][]float64
		sss       [][]string
		dss       [][]time.Time
		expectedO options
	}{
		{
			name:      "empty case",
			i:         []string{},
			o:         options{separator: '\t', scaleType: linear, chartType: pie},
			fss:       [][]float64{},
			sss:       [][]string{},
			dss:       [][]time.Time{},
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
			dss:       [][]time.Time{{}, {}, {}, {}},
			expectedO: options{separator: '\t', scaleType: linear, chartType: pie},
		},
	}

	for _, ts := range tests {
		fss, sss, dss, o := preprocess(ts.i, ts.o)
		if !reflect.DeepEqual(fss, ts.fss) {
			t.Errorf("'%v' failed: (floats) %v was not equal to %v", ts.name, fss, ts.fss)
		}
		if !reflect.DeepEqual(sss, ts.sss) {
			t.Errorf("'%v' failed: (strings) %v was not equal to %v", ts.name, sss, ts.sss)
		}
		if !reflect.DeepEqual(dss, ts.dss) {
			t.Errorf("'%v' failed: (times) %v was not equal to %v", ts.name, dss, ts.dss)
			if dss == nil {
				fmt.Println("thats' why")
			}
		}
		if !reflect.DeepEqual(o, ts.o) {
			t.Errorf("'%v' failed: %v was not equal to %v", ts.name, o, ts.expectedO)
		}
	}
}
