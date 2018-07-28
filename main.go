package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	o := mustResolveOptions(os.Args[1:])
	o, b := mustBuildChart(os.Stdin, o)
	tmpfile := mustNewTempFile()
	chartTempl := newChartTemplate(o.chartType)
	chartTempl.mustExecute(b, tmpfile)
	tmpfile.mustClose()
	tmpfile.mustRenameWithHTMLSuffix()
	mustOpen(tmpfile.url())
}

func mustOpen(url string) {
	if err := open.Run(url); err != nil {
		log.WithField("err", err).Fatalf("Could not open the default viewer; please configure open/xdg-open")
	}
}

func mustBuildChart(r io.Reader, o options) (options, bytes.Buffer) {
	_, o, b, err := buildChart(r, o)
	if err == nil && b.Len() == 0 {
		os.Exit(0)
	}
	if err != nil {
		log.Fatal(err)
	}
	return o, b
}

func buildChart(r io.Reader, o options) ([]string, options, bytes.Buffer, error) {
	d, o, lf, ls := preprocess(r, o)
	var b bytes.Buffer

	if o.debug {
		showDebug(ls, d, o, lf)
		return ls, o, b, nil
	}

	var err error
	var templ *template.Template
	var templData interface{}

	switch o.chartType {
	case pie:
		if len(d.fss) == 0 || (len(d.fss[0]) == 1 && len(d.sss) == 0 && len(d.tss) == 0) {
			return ls, o, b, fmt.Errorf("couldn't find values to plot")
		}
	case bar:
		if len(d.fss) == 0 || (len(d.fss[0]) == 1 && len(d.sss) == 0 && len(d.tss) == 0) {
			return ls, o, b, fmt.Errorf("couldn't find values to plot")
		}
	case line:
		if d.fss == nil || (d.sss == nil && d.tss == nil && len(d.fss[0]) < 2) {
			return ls, o, b, fmt.Errorf("couldn't find values to plot")
		}
	case scatter:
		if len(d.fss) == 0 {
			return ls, o, b, fmt.Errorf("couldn't find values to plot")
		}
	}

	templData, templ, err = cjsChart{inData{
		ChartType: o.chartType.string(),
		FSS:       d.fss,
		SSS:       d.sss,
		TSS:       d.tss,
		MinFSS:    d.minFSS,
		MaxFSS:    d.maxFSS,
		Title:     o.title,
		ScaleType: o.scaleType.string(),
		XLabel:    o.xLabel,
		YLabel:    o.yLabel,
		ZeroBased: o.zeroBased,
		ColorType: int(o.colorType),
	}}.chart()
	if err != nil {
		return ls, o, b, fmt.Errorf("couldn't construct chart because [%v]", err)
	}

	if err := templ.Execute(&b, templData); err != nil {
		return ls, o, b, fmt.Errorf("could't prepare ChartJS js code for chart: [%v]", err)
	}

	return ls, o, b, nil
}
