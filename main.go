package main

import (
	"bytes"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	var (
		opts           = mustResolveOptions(os.Args[1:])
		dataset, iOpts = preprocess(os.Stdin, opts)
		err            = assertChartable(dataset, iOpts)
	)
	if iOpts.debug || err != nil {
		showDebug(dataset, iOpts, err)
		os.Exit(0)
	}
	var (
		b          = mustBuildChart(dataset, iOpts)
		tmpfile    = mustNewTempFile()
		chartTempl = newChartTemplate(iOpts.chartType)
	)
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

func mustBuildChart(d dataset, o options) bytes.Buffer {
	b, err := buildChart(d, o)
	if err == nil && b.Len() == 0 {
		log.Println("Empty result; nothing to plot here.")
		os.Exit(0)
	}
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func buildChart(d dataset, o options) (bytes.Buffer, error) {
	templData, templ, err := cjsChart{inData{
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
		return bytes.Buffer{}, fmt.Errorf("couldn't construct chart because [%v]", err)
	}

	var b bytes.Buffer
	if err := templ.Execute(&b, templData); err != nil {
		return b, fmt.Errorf("could't prepare ChartJS js code for chart: [%v]", err)
	}

	return b, nil
}

func assertChartable(d dataset, opts options) error {
	switch opts.chartType {
	case pie:
		if len(d.fss) == 0 || (len(d.fss[0]) == 1 && len(d.sss) == 0 && len(d.tss) == 0) {
			return fmt.Errorf("couldn't find values to plot")
		}
	case bar:
		if len(d.fss) == 0 || (len(d.fss[0]) == 1 && len(d.sss) == 0 && len(d.tss) == 0) {
			return fmt.Errorf("couldn't find values to plot")
		}
	case line:
		if d.fss == nil || (d.sss == nil && d.tss == nil && len(d.fss[0]) < 2) {
			return fmt.Errorf("couldn't find values to plot")
		}
	case scatter:
		if len(d.fss) == 0 {
			return fmt.Errorf("couldn't find values to plot")
		}
	}
	return nil
}
