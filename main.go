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
	if err != nil {
		log.Fatal(err)
	}

	tmpfile, err := ioutil.TempFile("", "chartData")
	if err != nil {
		log.WithField("err", err).Fatalf("Could not create temporary file to store the chart.")
	}

	if err = baseTemplate.Execute(tmpfile, b.String()); err != nil {
		log.WithField("err", err).Fatalf("Could not write chart to temporary file.")
	}
	if err = tmpfile.Close(); err != nil {
		log.WithField("err", err).Fatalf("Could not close temporary file after saving chart to it.")
	}

	open.Run("file://" + tmpfile.Name())
}

func buildChart(i []string, o options) (bytes.Buffer, error) {
	fss, sss, tss, minFSS, maxFSS, o := preprocess(i, o)

	var b bytes.Buffer
	var err error
	var templ *template.Template
	var templData interface{}

	switch o.chartType {
	case pie:
		templData, templ, err = setupPie(fss, sss, o.title)
	case bar:
		templData, templ, err = setupBar(fss, sss, o.title, o.scaleType, o.xLabel, o.yLabel)
	case line:
		templData, templ, err = setupLine(fss, sss, tss, o.title, o.scaleType, o.xLabel, o.yLabel)
	case scatter:
		templData, templ, err = setupScatter(fss, sss, tss, minFSS, maxFSS, o.title, o.scaleType, o.xLabel, o.yLabel)
	}
	if err != nil {
		return b, fmt.Errorf("Could not construct chart because [%v]", err)
	}
	if err := templ.Execute(&b, templData); err != nil {
		return b, fmt.Errorf("Could not prepare ChartJS js code for chart: [%v]", err)
	}

	return b, nil
}
