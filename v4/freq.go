package main

import (
	"sort"

	"github.com/marianogappa/chart/v4/format"
)

func preprocessFreq(isss [][]string, lineFormat format.LineFormat) ([][]float64, [][]string, format.LineFormat) {
	fss := [][]float64{}
	sss := [][]string{}

	is := make(map[string]int)
	fqs := freqs{}
	for _, ss := range isss {
		if len(ss) == 0 {
			break //TODO this probably shouldn't happen
		}
		if _, ok := is[ss[0]]; !ok {
			is[ss[0]] = len(fqs.fs)
			fqs.fs = append(fqs.fs, freq{s: ss[0], f: 1})
		} else {
			fqs.fs[is[ss[0]]].f++
		}
	}

	sort.Sort(fqs)

	for i := 0; i < 9; i++ {
		if i >= len(fqs.fs) {
			break
		}
		fss = append(fss, []float64{fqs.fs[i].f})
		sss = append(sss, []string{fqs.fs[i].s})
	}

	if len(fqs.fs) > 10 {
		sum := float64(0)
		for i := 9; i < len(fqs.fs); i++ {
			sum += fqs.fs[i].f
		}
		fss = append(fss, []float64{sum})
		sss = append(sss, []string{"Other"})
	} else if len(fqs.fs) == 10 {
		fss = append(fss, []float64{fqs.fs[9].f})
		sss = append(sss, []string{fqs.fs[9].s})
	}

	// Updates lineFormat
	lineFormat.ColTypes = append(lineFormat.ColTypes, format.Float)
	lineFormat.FloatCount++
	lineFormat.HasFloats = true

	return fss, sss, lineFormat
}

type freq struct {
	s string
	f float64
}

type freqs struct {
	fs []freq
}

func (f freqs) Len() int      { return len(f.fs) }
func (f freqs) Swap(i, j int) { f.fs[i], f.fs[j] = f.fs[j], f.fs[i] }

func (f freqs) Less(i, j int) bool {
	if f.fs[i].f == f.fs[j].f {
		return i < j
	}
	return f.fs[i].f > f.fs[j].f
}
