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
			dss:       [][]time.Time{{}, {}, {}, {}}, //TODO fix the frequence case
			expectedO: options{separator: '\t', scaleType: linear, chartType: pie},
		},
		{
			name: "dates and floats with many spaces in between",
			i: []string{"2016-08-29	0.0125", "2016-09-06	0.0272", "2016-09-07	0.0000", "2016-09-08	0.0000"},
			o: options{separator: '\t', scaleType: linear, chartType: pie, dateFormat: "2006-01-02"},
			fss: [][]float64{
				[]float64{0.0125}, []float64{0.0272}, []float64{0}, []float64{0},
			},
			sss: [][]string{
				{}, {}, {}, {},
			},
			dss: [][]time.Time{
				{tp("2006-01-02", "2016-08-29")}, {tp("2006-01-02", "2016-09-06")}, {tp("2006-01-02", "2016-09-07")}, {tp("2006-01-02", "2016-09-08")},
			},
			expectedO: options{separator: '\t', scaleType: linear, chartType: pie, dateFormat: "2006-01-02"},
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
		if !reflect.DeepEqual(o, ts.expectedO) {
			t.Errorf("'%v' failed: %v was not equal to %v", ts.name, o, ts.expectedO)
		}
	}
}

func tp(tf, t string) time.Time {
	rt, _ := time.Parse(tf, t)
	return rt
}
