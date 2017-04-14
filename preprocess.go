package main

import (
	"strings"
	"time"
)

func preprocess(i []string, o options) ([][]float64, [][]string, [][]time.Time, []float64, []float64, options, string) {
	var fss [][]float64
	var sss [][]string
	var tss [][]time.Time
	var minFSS []float64
	var maxFSS []float64

	sep := o.separator
	lf := parseFormat(i, sep, o.dateFormat)
	for _, l := range i {
		fs, ss, ts, err := parseLine(l, lf, sep, o.dateFormat)
		if err != nil {
			continue
		}

		for i, f := range fs {
			if len(minFSS) == i {
				minFSS = append(minFSS, f)
			}
			if len(maxFSS) == i {
				maxFSS = append(maxFSS, f)
			}
			if f < minFSS[i] {
				minFSS[i] = f
			}
			if f > maxFSS[i] {
				maxFSS[i] = f
			}
		}

		fss = append(fss, fs)
		sss = append(sss, ss)
		tss = append(tss, ts)
	}
	o.chartType = resolveChartType(o.chartType, lf, fss, sss)

	if o.chartType == bar {
		o.zeroBased = true // https://github.com/marianogappa/chart/issues/11
	}

	if strings.Index(lf, "f") == -1 {
		fss, sss = preprocessFreq(sss)
	}

	nilSSS := true
	for _, ss := range sss {
		if len(ss) > 0 {
			nilSSS = false
			break
		}
	}
	if nilSSS {
		sss = nil
	}

	nilFSS := true
	for _, fs := range fss {
		if len(fs) > 0 {
			nilFSS = false
			break
		}
	}
	if nilFSS {
		fss = nil
	}

	nilTSS := true
	for _, ts := range tss {
		if len(ts) > 0 {
			nilTSS = false
			break
		}
	}
	if nilTSS {
		tss = nil
	}

	return fss, sss, tss, minFSS, maxFSS, o, lf
}
