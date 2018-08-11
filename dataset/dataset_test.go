package dataset

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/marianogappa/chart/format"
)

func TestDataset(t *testing.T) {
	tests := []struct {
		name          string
		i             string
		rawLineFormat string
		sep           rune
		dateFormat    string
		fss           [][]float64
		sss           [][]string
		tss           [][]time.Time
		minFSS        []float64
		maxFSS        []float64
	}{
		{
			name: "format: ff; line count: 3; separator: tab",
			i: `
			1	6
			3	4
			5	2
			`,
			rawLineFormat: "ff",
			sep:           '\t',
			dateFormat:    "",
			fss:           [][]float64{{1, 6}, {3, 4}, {5, 2}},
			sss:           nil,
			tss:           nil,
			minFSS:        []float64{1, 2},
			maxFSS:        []float64{5, 6},
		},
		{
			name: "format: df; line count: 3; separator: tab",
			i: `
			2016-08-29	4
			2016-09-08	4
			2016-09-06	3
			2016-09-07	3
			`,
			rawLineFormat: "df",
			sep:           '\t',
			dateFormat:    "2006-01-02",
			fss:           [][]float64{{4}, {4}, {3}, {3}},
			sss:           nil,
			tss: [][]time.Time{
				{tp("2006-01-02", "2016-08-29")},
				{tp("2006-01-02", "2016-09-08")},
				{tp("2006-01-02", "2016-09-06")},
				{tp("2006-01-02", "2016-09-07")}},
			minFSS: []float64{3},
			maxFSS: []float64{4},
		},
		{
			name: "just strings",
			i: `a
		b
		c
		a
		`,
			rawLineFormat: "s",
			sep:           '\t',
			dateFormat:    "",
			sss:           [][]string{{"a"}, {"b"}, {"c"}, {"a"}},
			tss:           nil,
			minFSS:        nil,
			maxFSS:        nil,
		},
		{
			name: "dates and floats with many spaces in between",
			i: `2016-08-29	0.0125
		2016-09-06	0.0272
		2016-09-07	0.0000
		2016-09-08	0.0000
		`,
			rawLineFormat: "df",
			sep:           '\t',
			dateFormat:    "2006-01-02",
			fss: [][]float64{
				{0.0125}, {0.0272}, {0}, {0},
			},
			sss: nil,
			tss: [][]time.Time{
				{tp("2006-01-02", "2016-08-29")}, {tp("2006-01-02", "2016-09-06")}, {tp("2006-01-02", "2016-09-07")}, {tp("2006-01-02", "2016-09-08")},
			},
			minFSS: []float64{0.0000},
			maxFSS: []float64{0.0272},
		},
	}

	for _, ts := range tests {
		lf, _ := format.NewLineFormat(ts.rawLineFormat, ts.sep, ts.dateFormat)
		d, err := newDataset(bytes.NewReader([]byte(ts.i)), lf)

		if err != nil {
			t.Errorf("'%v' failed: error reading dataset %v", ts.name, err)
		}
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
	}
}

func tp(tf, t string) time.Time {
	rt, _ := time.Parse(tf, t)
	return rt
}
