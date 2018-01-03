package main

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestPreprocess(t *testing.T) {
	tests := []struct {
		name       string
		i          string
		o          options
		fss        [][]float64
		sss        [][]string
		tss        [][]time.Time
		minFSS     []float64
		maxFSS     []float64
		expectedO  options
		expectedLF string
	}{
		{
			name:      "empty case",
			i:         ``,
			o:         options{separator: '\t', scaleType: linear, chartType: pie},
			fss:       nil,
			sss:       nil,
			tss:       nil,
			minFSS:    nil,
			maxFSS:    nil,
			expectedO: options{separator: '\t', scaleType: linear, chartType: pie},
		},
		{
			name: "sorts scatter line charts",
			i: `2	4
1	3
4	4
3	3
`,
			o:          options{separator: '\t', scaleType: linear, chartType: line},
			fss:        [][]float64{{1, 3}, {2, 4}, {3, 3}, {4, 4}},
			sss:        nil,
			tss:        nil,
			minFSS:     []float64{1, 3},
			maxFSS:     []float64{4, 4},
			expectedO:  options{separator: '\t', scaleType: linear, chartType: line},
			expectedLF: "ff",
		},
		{
			name: "sorts scatter line charts with time series",
			i: `2016-08-29	4
2016-09-08	4
2016-09-06	3
2016-09-07	3
`,
			o:   options{separator: '\t', scaleType: linear, chartType: line, dateFormat: "2006-01-02"},
			fss: [][]float64{{4}, {3}, {3}, {4}},
			sss: nil,
			tss: [][]time.Time{
				{tp("2006-01-02", "2016-08-29")}, {tp("2006-01-02", "2016-09-06")}, {tp("2006-01-02", "2016-09-07")}, {tp("2006-01-02", "2016-09-08")},
			},
			minFSS:     []float64{3},
			maxFSS:     []float64{4},
			expectedO:  options{separator: '\t', scaleType: linear, chartType: line, dateFormat: "2006-01-02"},
			expectedLF: "df",
		},
		{
			name: "just strings",
			i: `a
		b
		c
		a
		`,
			o: options{separator: '\t', scaleType: linear, chartType: pie},
			fss: [][]float64{
				{2}, {1}, {1},
			},
			sss:        [][]string{{"a"}, {"b"}, {"c"}},
			tss:        nil,
			minFSS:     nil,
			maxFSS:     nil,
			expectedO:  options{separator: '\t', scaleType: linear, chartType: pie},
			expectedLF: "s",
		},
		{
			name: "dates and floats with many spaces in between",
			i: `2016-08-29	0.0125
		2016-09-06	0.0272
		2016-09-07	0.0000
		2016-09-08	0.0000
		`,
			o: options{separator: '\t', scaleType: linear, chartType: pie, dateFormat: "2006-01-02"},
			fss: [][]float64{
				{0.0125}, {0.0272}, {0}, {0},
			},
			sss: nil,
			tss: [][]time.Time{
				{tp("2006-01-02", "2016-08-29")}, {tp("2006-01-02", "2016-09-06")}, {tp("2006-01-02", "2016-09-07")}, {tp("2006-01-02", "2016-09-08")},
			},
			minFSS:     []float64{0.0000},
			maxFSS:     []float64{0.0272},
			expectedO:  options{separator: '\t', scaleType: linear, chartType: pie, dateFormat: "2006-01-02"},
			expectedLF: "df",
		},
	}

	for _, ts := range tests {
		d, o, lf, _ := preprocess(strings.NewReader(ts.i), ts.o)

		if !reflect.DeepEqual(d.fss, ts.fss) {
			t.Errorf("'%v' failed: (floats) %v was not equal to %v", ts.name, d.fss, ts.fss)
		}
		if !reflect.DeepEqual(d.sss, ts.sss) {
			t.Errorf("'%v' failed: (strings) %v was not equal to %v", ts.name, d.sss, ts.sss)
		}
		if !reflect.DeepEqual(d.tss, ts.tss) {
			t.Errorf("'%v' failed: (times) %v was not equal to %v", ts.name, d.tss, ts.tss)
		}
		if !reflect.DeepEqual(d.minFSS, ts.minFSS) {
			t.Errorf("'%v' failed: (min floats) %v was not equal to %v", ts.name, d.minFSS, ts.minFSS)
		}
		if !reflect.DeepEqual(d.maxFSS, ts.maxFSS) {
			t.Errorf("'%v' failed: (max floats) %v was not equal to %v", ts.name, d.maxFSS, ts.maxFSS)
		}
		if !reflect.DeepEqual(o, ts.expectedO) {
			t.Errorf("'%v' failed: %v was not equal to %v", ts.name, o, ts.expectedO)
		}
		if !reflect.DeepEqual(lf, ts.expectedLF) {
			t.Errorf("'%v' failed: %v was not equal to %v", ts.name, lf, ts.expectedLF)
		}
	}
}

func tp(tf, t string) time.Time {
	rt, _ := time.Parse(tf, t)
	return rt
}
