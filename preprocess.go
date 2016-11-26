package main

import (
	"strings"
	"time"
)

func preprocess(i []string, o options) ([][]float64, [][]string, [][]time.Time, options) {
	var fss [][]float64
	var sss [][]string
	dss := [][]time.Time{}

	sep := o.separator
	lf := parseFormat(i, sep, o.dateFormat)
	for _, l := range i {
		fs, ss, ds, err := parseLine(l, lf, sep, o.dateFormat)
		if err != nil {
			break
		}
		fss = append(fss, fs)
		sss = append(sss, ss)
		dss = append(dss, ds)
	}

	o.chartType = resolveChartType(o.chartType, lf, fss, sss)

	if strings.Index(lf, "f") == -1 {
		fss, sss = preprocessFreq(sss)
	}

	return fss, sss, dss, o
}
