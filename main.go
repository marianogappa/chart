package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/skratchdot/open-golang/open"
)

func main() {
	o := mustResolveOptions(os.Args[1:])
	i := readInput(os.Stdin)
	fss, sss, o := preprocess(i, o)

	var err error
	var templ *template.Template
	var templData interface{}

	switch o.chartType {
	case pie:
		templData, templ, err = setupPie(fss, sss, o.title)
	case bar:
		templData, templ, err = setupBar(fss, sss, o.title, o.scaleType)
	}
	if err != nil {
		log.Fatalf("Could not construct chart because [%v]", err)
	}
	var b bytes.Buffer
	if err := templ.Execute(&b, templData); err != nil {
		log.Fatalf("Could not prepare ChartJS js code for chart: [%v]", err)
	}

	tmpfile, err := ioutil.TempFile("", "chartData")
	if err != nil {
		log.Fatalf("Could not create temporary file to store the chart: [%v]", err)
	}

	// TODO is it worth to introduce delay and race condition to delete a temp file?
	// defer func() { time.Sleep(5 * time.Second); os.Remove(tmpfile.Name()) }()

	if err = baseTemplate.Execute(tmpfile, b.String()); err != nil {
		log.Fatalf("Could not write chart to temporary file: [%v]", err)
	}
	if err = tmpfile.Close(); err != nil {
		log.Fatalf("Could not close temporary file after saving chart to it: [%v]", err)
	}

	open.Run("file://" + tmpfile.Name())
}
