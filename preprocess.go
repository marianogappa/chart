package main

import (
	"io"
	"sort"
	"strings"
	"time"
)

type dataset struct {
	fss    [][]float64
	sss    [][]string
	tss    [][]time.Time
	minFSS []float64
	maxFSS []float64
}

func newDataset() *dataset {
	return &dataset{
		fss:    make([][]float64, 0, 500),
		sss:    make([][]string, 0, 500),
		tss:    make([][]time.Time, 0, 500),
		minFSS: make([]float64, 0, 500),
		maxFSS: make([]float64, 0, 500),
	}
}

func (d *dataset) Len() int {
	if d.fss == nil {
		return len(d.tss)
	}
	return len(d.fss)
}

func (d *dataset) Less(i, j int) bool {
	if d.tss == nil {
		return d.fss[i][0] < d.fss[j][0]
	}
	return d.tss[i][0].Before(d.tss[j][0])
}

func (d *dataset) Swap(i, j int) {
	if d.fss != nil {
		d.fss[i], d.fss[j] = d.fss[j], d.fss[i]
	}
	if d.tss != nil {
		d.tss[i], d.tss[j] = d.tss[j], d.tss[i]
	}
	if d.sss != nil {
		d.sss[i], d.sss[j] = d.sss[j], d.sss[i]
	}
}

func (d *dataset) hasFloats() bool  { return len(d.fss) > 0 }
func (d *dataset) hasStrings() bool { return len(d.sss) > 0 }
func (d *dataset) hasTimes() bool   { return len(d.tss) > 0 }
func (d *dataset) timeFieldLen() int {
	if !d.hasTimes() {
		return 0
	}
	return len(d.tss[0])
}
func (d *dataset) floatFieldLen() int {
	if !d.hasFloats() {
		return 0
	}
	return len(d.fss[0])
}
func (d *dataset) canBeScatterLine() bool {
	return d.floatFieldLen()+d.timeFieldLen() >= 2
}

func preprocess(r io.Reader, o options) (dataset, options, string, []string) {
	var (
		d                      = newDataset()
		sep                    = o.separator
		ls, lf                 = readAndParseFormat(r, sep, o.dateFormat)
		nilSSS, nilFSS, nilTSS = true, true, true
	)

	for _, l := range ls {
		fs, ss, ts, err := parseLine(l, lf, sep, o.dateFormat)
		if err != nil {
			continue
		}
		if nilSSS && len(ss) > 0 {
			nilSSS = false
		}
		if nilFSS && len(fs) > 0 {
			nilFSS = false
		}
		if nilTSS && len(ts) > 0 {
			nilTSS = false
		}

		for i, f := range fs {
			if len(d.minFSS) == i {
				d.minFSS = append(d.minFSS, f)
			}
			if len(d.maxFSS) == i {
				d.maxFSS = append(d.maxFSS, f)
			}
			if f < d.minFSS[i] {
				d.minFSS[i] = f
			}
			if f > d.maxFSS[i] {
				d.maxFSS[i] = f
			}
		}

		d.fss = append(d.fss, fs)
		d.sss = append(d.sss, ss)
		d.tss = append(d.tss, ts)
	}
	if nilSSS {
		d.sss = nil
	}
	if nilFSS {
		d.fss = nil
		d.minFSS = nil
		d.maxFSS = nil
	}
	if nilTSS {
		d.tss = nil
	}

	o.chartType = resolveChartType(o.chartType, lf, d.fss, d.sss)

	if o.chartType == bar {
		o.zeroBased = true // https://github.com/marianogappa/chart/issues/11
	}

	if strings.Index(lf, "f") == -1 && len(d.sss) > 0 {
		d.fss, d.sss = preprocessFreq(d.sss)
	}

	if o.chartType == line && d.canBeScatterLine() {
		sort.Sort(d)
	}

	return *d, o, lf, ls
}
