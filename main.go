package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/skratchdot/open-golang/open"
)

func main() {
	chartType := mustParseChartType()
	title, displayTitle := mustParseTitle()
	// defer func() { time.Sleep(5 * time.Second); os.Remove(tmpfile.Name()) }()

	var finalString string
	var err error
	switch chartType {
	case "pie":
		finalString, err = setupPie(title, displayTitle, '\t')
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

func mustParseChartType() string {
	if len(os.Args) < 2 {
		return "pie"
	}

	switch os.Args[1] {
	case "pie":
		return "pie"
	default:
		return "pie"
	}
}

func mustParseTitle() (string, bool) {
	if len(os.Args) < 3 {
		return "", false
	}
	return os.Args[2], true
}
