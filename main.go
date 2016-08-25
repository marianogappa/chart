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

	// defer func() { time.Sleep(5 * time.Second); os.Remove(tmpfile.Name()) }()

	var finalString string
	var err error
	switch o.chartType {
	case pie:
		finalString, err = setupPie(i, o.title, len(o.title) > 0, rune(o.separator[0]), o.invert)
	case bar:
		finalString, err = setupBar(i, o.title, len(o.title) > 0, rune(o.separator[0]), o.scaleType, o.invert)
	}
	if err != nil {
		log.Fatal(err)
	}

	tmpfile, err := ioutil.TempFile("", "chartData")
	if err != nil {
		log.Fatal(err)
	}

	if err = baseTemplate.Execute(tmpfile, finalString); err != nil {
		log.Fatal(err)
	}
	if err = tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	open.Run("file://" + tmpfile.Name())
}
