package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/skratchdot/open-golang/open"
)

func main() {
	o := mustResolveOptions(os.Args[1:])
	i := mustReadInput(os.Stdin)
	fss, sss, o := preprocess(i, o)

	// defer func() { time.Sleep(5 * time.Second); os.Remove(tmpfile.Name()) }()

	var html string
	var err error
	switch o.chartType {
	case pie:
		html, err = setupPie(fss, sss, o.title, len(o.title) > 0)
	case bar:
		html, err = setupBar(fss, sss, o.title, len(o.title) > 0, o.scaleType)
	}
	if err != nil {
		log.Fatal(err)
	}

	tmpfile, err := ioutil.TempFile("", "chartData")
	if err != nil {
		log.Fatal(err)
	}

	if err = baseTemplate.Execute(tmpfile, html); err != nil {
		log.Fatal(err)
	}
	if err = tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	open.Run("file://" + tmpfile.Name())
}
