package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	o := mustResolveOptions(os.Args[1:])
	i := readInput(os.Stdin)

	b, err := buildChart(i, o)
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

	if err = baseTemplate.Execute(tmpfile, b.String()); err != nil {
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

func buildChart(i []string, o options) (bytes.Buffer, error) {
	fss, sss, tss, minFSS, maxFSS, o, lf := preprocess(i, o)
	var b bytes.Buffer

	if o.debug {
		showDebug(i, fss, sss, tss, minFSS, maxFSS, o, lf)
		return b, nil
	}

	var err error
	var templ *template.Template
	var templData interface{}

	switch o.chartType {
	case pie:
		if len(fss) == 0 || (len(fss[0]) == 1 && len(sss) == 0 && len(tss) == 0) {
			return b, fmt.Errorf("couldn't find values to plot")
		}
	case bar:
		if len(fss) == 0 || (len(fss[0]) == 1 && len(sss) == 0 && len(tss) == 0) {
			return b, fmt.Errorf("couldn't find values to plot")
		}
	case line:
		if fss == nil || (sss == nil && tss == nil && len(fss[0]) < 2) {
			return b, fmt.Errorf("couldn't find values to plot")
		}
	case scatter:
		if len(fss) == 0 {
			return b, fmt.Errorf("couldn't find values to plot")
		}
	}
	if err != nil {
		return b, fmt.Errorf("could not construct chart because [%v]", err)
	}

	templData, templ, err = cjsChart{inData{
		ChartType: o.chartType.string(),
		FSS:       fss,
		SSS:       sss,
		TSS:       tss,
		MinFSS:    minFSS,
		MaxFSS:    maxFSS,
		Title:     o.title,
		ScaleType: o.scaleType.string(),
		XLabel:    o.xLabel,
		YLabel:    o.yLabel,
		ZeroBased: o.zeroBased,
	}}.chart()

	if err := templ.Execute(&b, templData); err != nil {
		return b, fmt.Errorf("could not prepare ChartJS js code for chart: [%v]", err)
	}

	return b, nil
}
