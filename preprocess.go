package main

import "strings"

func preprocess(i []string, o options) ([][]float64, [][]string, options) {
	var fss [][]float64
	var sss [][]string

	sep := o.separator
	lf := parseFormat(i, sep)
	for _, l := range i {
		fs, ss, err := parseLine(l, lf, sep)
		if err != nil {
			break
		}
		fss = append(fss, fs)
		sss = append(sss, ss)
	}

	o.chartType = resolveChartType(o.chartType, lf, fss, sss)

	if strings.Index(lf, "f") == -1 {
		fss, sss = preprocessFreq(sss)
	}

	return fss, sss, o
}
