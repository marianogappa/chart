package main

import (
	"fmt"
)

func showDebug(d dataset, o options, err error) {
	if err != nil {
		fmt.Printf("Error trying to chart: %v\n", err)
	}
	fcn, scn, tcn, rn := 0, 0, 0, 0
	if len(d.fss) > 0 {
		rn = len(d.fss)
		fcn = len(d.fss[0])
	}
	if len(d.sss) > 0 {
		rn = len(d.sss)
		scn = len(d.sss[0])
	}
	if len(d.tss) > 0 {
		rn = len(d.tss)
		tcn = len(d.tss[0])
	}
	fmt.Printf("Lines read\t%v\n", d.stdinLen)
	fmt.Printf("Line format inferred\t%v\n", d.lf)
	fmt.Printf("Lines used\t%v\n", rn)
	fmt.Printf("Float column count\t%v\n", fcn)
	fmt.Printf("String column count\t%v\n", scn)
	fmt.Printf("Date/Time column count\t%v\n", tcn)

	if o.title != "" {
		fmt.Printf("Chart title\t%v\n", o.title)
	}
	if o.xLabel != "" {
		fmt.Printf("Chart horizontal axis label\t%v\n", o.xLabel)
	}
	if o.yLabel != "" {
		fmt.Printf("Chart vertical axis label\t%v\n", o.yLabel)
	}
	if o.dateFormat != "" {
		fmt.Printf("Date format\t%v\n", o.dateFormat)
	}
	switch o.chartType {
	case pie:
		fmt.Println("Chart type\tpie")
	case bar:
		fmt.Println("Chart type\tbar")
	case line:
		fmt.Println("Chart type\tline")
	case scatter:
		fmt.Println("Chart type\tscatter")
	default:
		fallthrough
	case undefinedChartType:
		fmt.Println("Chart type\t???")
	}
	switch o.scaleType {
	case linear:
		fmt.Println("Scale type\tlinear")
	case logarithmic:
		fmt.Println("Scale type\tlogarithmic")
	}
	switch o.separator {
	case '\t':
		fmt.Printf("Separator\t[tab]\n")
	case ' ':
		fmt.Printf("Separator\t[space]\n")
	case ',':
		fmt.Printf("Separator\t[comma]\n")
	case ';':
		fmt.Printf("Separator\t[semicolon]\n")
	}
}
