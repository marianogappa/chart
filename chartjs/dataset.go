package chartjs

import "time"

type dataset struct {
	ChartType string
	FSS       [][]float64
	TSS       [][]time.Time
	SSS       [][]string
	MinFSS    []float64
	MaxFSS    []float64
	Title     string
	ScaleType string
	XLabel    string
	YLabel    string
	ZeroBased bool
	ColorType int
}

func (d dataset) Len() int {
	if d.FSS == nil {
		return len(d.TSS)
	}
	return len(d.FSS)
}

func (d dataset) Less(i, j int) bool {
	if d.TSS == nil {
		return d.FSS[i][0] < d.FSS[j][0]
	}
	return d.TSS[i][0].Before(d.TSS[j][0])
}

func (d dataset) Swap(i, j int) {
	if d.FSS != nil {
		d.FSS[i], d.FSS[j] = d.FSS[j], d.FSS[i]
	}
	if d.TSS != nil {
		d.TSS[i], d.TSS[j] = d.TSS[j], d.TSS[i]
	}
	if d.SSS != nil {
		d.SSS[i], d.SSS[j] = d.SSS[j], d.SSS[i]
	}
}

func (d dataset) hasFloats() bool  { return len(d.FSS) > 0 }
func (d dataset) hasStrings() bool { return len(d.SSS) > 0 }
func (d dataset) hasTimes() bool   { return len(d.TSS) > 0 }
func (d dataset) timeFieldLen() int {
	if !d.hasTimes() {
		return 0
	}
	return len(d.TSS[0])
}
func (d dataset) floatFieldLen() int {
	if !d.hasFloats() {
		return 0
	}
	return len(d.FSS[0])
}
func (d dataset) canBeScatterLine() bool {
	return d.floatFieldLen()+d.timeFieldLen() >= 2
}
