package main

import (
	"bytes"
	"fmt"

	"github.com/marianogappa/chart/v4/dataset"
)

func renderDebug(d dataset.Dataset, o options, err error) string {
	var buffer bytes.Buffer
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Error trying to chart: %v\n", err))
	}
	fcn, scn, tcn, rn := 0, 0, 0, 0
	if len(d.FSS) > 0 {
		rn = len(d.FSS)
		fcn = len(d.FSS[0])
	}
	if len(d.SSS) > 0 {
		rn = len(d.SSS)
		scn = len(d.SSS[0])
	}
	if len(d.TSS) > 0 {
		rn = len(d.TSS)
		tcn = len(d.TSS[0])
	}
	buffer.WriteString(fmt.Sprintf("Lines read\t%v\n", d.StdinLen))
	buffer.WriteString(fmt.Sprintf("Line format inferred\t%v\n", o.lineFormat.String()))
	buffer.WriteString(fmt.Sprintf("Lines used\t%v\n", rn))
	buffer.WriteString(fmt.Sprintf("Float column count\t%v\n", fcn))
	buffer.WriteString(fmt.Sprintf("String column count\t%v\n", scn))
	buffer.WriteString(fmt.Sprintf("Date/Time column count\t%v\n", tcn))

	if o.title != "" {
		buffer.WriteString(fmt.Sprintf("Chart title\t%v\n", o.title))
	}
	if o.xLabel != "" {
		buffer.WriteString(fmt.Sprintf("Chart horizontal axis label\t%v\n", o.xLabel))
	}
	if o.yLabel != "" {
		buffer.WriteString(fmt.Sprintf("Chart vertical axis label\t%v\n", o.yLabel))
	}
	if o.dateFormat != "" {
		buffer.WriteString(fmt.Sprintf("Date format\t%v\n", o.dateFormat))
	}
	switch o.chartType {
	case pie:
		buffer.WriteString(fmt.Sprintf("Chart type\tpie\n"))
	case bar:
		buffer.WriteString(fmt.Sprintf("Chart type\tbar\n"))
	case line:
		buffer.WriteString(fmt.Sprintf("Chart type\tline\n"))
	case scatter:
		buffer.WriteString(fmt.Sprintf("Chart type\tscatter\n"))
	case undefinedChartType:
		buffer.WriteString(fmt.Sprintf("Chart type\t???\n"))
	}
	switch o.scaleType {
	case linear:
		buffer.WriteString(fmt.Sprintf("Scale type\tlinear\n"))
	case logarithmic:
		buffer.WriteString(fmt.Sprintf("Scale type\tlogarithmic\n"))
	}
	switch o.separator {
	case '\t':
		buffer.WriteString(fmt.Sprintf("Separator\t[tab]\n"))
	case ' ':
		buffer.WriteString(fmt.Sprintf("Separator\t[space]\n"))
	case ',':
		buffer.WriteString(fmt.Sprintf("Separator\t[comma]\n"))
	case ';':
		buffer.WriteString(fmt.Sprintf("Separator\t[semicolon]\n"))
	}

	return buffer.String()
}
