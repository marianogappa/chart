package dataset

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/marianogappa/chart/v4/format"
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
		},
		{
			name: "just strings",
			i: `a
		b
		c
		a`,
			rawLineFormat: "s",
			sep:           '\t',
			dateFormat:    "",
			sss:           [][]string{{"a"}, {"b"}, {"c"}, {"a"}},
			tss:           nil,
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
		},
	}

	for _, ts := range tests {
		lf, _ := format.NewLineFormat(ts.rawLineFormat, ts.sep, ts.dateFormat)
		d, err := New(bytes.NewReader([]byte(ts.i)), lf)

		if err != nil {
			t.Errorf("'%v' failed: error reading dataset %v", ts.name, err)
		}
		if !reflect.DeepEqual(d.FSS, ts.fss) {
			t.Errorf("'%v' failed: (floats) %v was not equal to %v", ts.name, d.FSS, ts.fss)
		}
		if !reflect.DeepEqual(d.SSS, ts.sss) {
			t.Errorf("'%v' failed: (strings) %v was not equal to %v", ts.name, d.SSS, ts.sss)
		}
		if !reflect.DeepEqual(d.TSS, ts.tss) {
			t.Errorf("'%v' failed: (times) %v was not equal to %v", ts.name, d.TSS, ts.tss)
		}
	}
}

func tp(tf, t string) time.Time {
	rt, _ := time.Parse(tf, t)
	return rt
}
