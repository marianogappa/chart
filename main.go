package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	o := mustResolveOptions(os.Args[1:])

	_, o, b, err := buildChart(os.Stdin, o)
	if err == nil && b.Len() == 0 {
		os.Exit(0)
	}
	if err != nil {
		log.Fatal(err)
	}

	tmpfile, err := ioutil.TempFile("", "chartData")
	if err != nil {
		log.WithField("err", err).Fatalf("Could not create temporary file to store the chart.")
	}

	if _, err := tmpfile.WriteString(baseTemplateHeaderString); err != nil {
		log.WithField("err", err).Fatalf("Could not write header to temporary file.")
	}

	t := baseTemplate
	if o.chartType == pie {
		t = basePieTemplate
	}
	if err = t.Execute(tmpfile, b.String()); err != nil {
		log.WithField("err", err).Fatalf("Could not write chart to temporary file.")
	}

	if _, err := tmpfile.WriteString(baseTemplateFooterString); err != nil {
		log.WithField("err", err).Fatalf("Could not write footer to temporary file.")
	}

	if err = tmpfile.Close(); err != nil {
		log.WithField("err", err).Fatalf("Could not close temporary file after saving chart to it.")
	}

	newName := tmpfile.Name() + ".html"
	if err = os.Rename(tmpfile.Name(), newName); err != nil {
		log.WithField("err", err).Fatalf("Could not add html extension to the temporary file.")
	}

	open.Run("file://" + newName)
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
	if err != nil {
		return ls, o, b, fmt.Errorf("could not construct chart because [%v]", err)
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

	if err := templ.Execute(&b, templData); err != nil {
		return ls, o, b, fmt.Errorf("could not prepare ChartJS js code for chart: [%v]", err)
	}

	return ls, o, b, nil
}
