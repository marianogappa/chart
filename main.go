package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/skratchdot/open-golang/open"
)

type options struct {
	title       string
	logarithmic bool
	typ         string
}

func main() {
	titleHelp := "Sets the title for the chart."
	logarithmicHelp := "Sets logarithmic scale for the y-axis."
	typHelp := "Sets the chart type; default is 'pie'; can be pie/bar."

	var o options

	fs := flag.NewFlagSet("params", flag.ContinueOnError)
	fs.StringVar(&o.title, "title", o.title, titleHelp)
	fs.StringVar(&o.title, "t", o.title, titleHelp)
	fs.BoolVar(&o.logarithmic, "log", o.logarithmic, logarithmicHelp)
	fs.BoolVar(&o.logarithmic, "l", o.logarithmic, logarithmicHelp)
	fs.StringVar(&o.typ, "type", o.typ, typHelp)
	fs.StringVar(&o.typ, "y", o.typ, typHelp)

	err := fs.Parse(fromFirstDash(os.Args[1:]))
	if err != nil {
		log.Fatal(err)
	}

	args := fs.Args()
	args = append(args, untilFirstDash(os.Args[1:])...)
	chartType := mustParseChartType(o.typ, args)
	title, displayTitle := mustParseTitle(o.title)
	scale := mustParseScale(o.logarithmic, args)
	// defer func() { time.Sleep(5 * time.Second); os.Remove(tmpfile.Name()) }()

	var finalString string
	switch chartType {
	case "pie":
		finalString, err = setupPie(title, displayTitle, '\t')
	case "bar":
		finalString, err = setupBar(title, displayTitle, '\t', scale)
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

func mustParseChartType(typ string, args []string) string {
	switch typ {
	case "pie":
		return "pie"
	case "bar":
		return "bar"
	default:
		if isInSlice("pie", args) {
			return "pie"
		} else if isInSlice("bar", args) {
			return "bar"
		}
	}
	return "pie"
}

func mustParseTitle(title string) (string, bool) {
	if len(title) == 0 {
		return "", false
	}
	return title, true
}

func mustParseScale(log bool, args []string) string {
	if log || isInSlice("log", args) {
		return "logarithmic"
	}

	return "linear"
}

func isInSlice(s string, ss []string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

// flag package doesn't read flags if first argument is not a flag :(
func fromFirstDash(as []string) []string {
	for i, v := range as {
		if v[0] == '-' {
			return as[i:]
		}
	}
	return []string{}
}

// flag package doesn't know about args before first dash :(
func untilFirstDash(as []string) []string {
	for i, v := range as {
		if v[0] == '-' {
			return as[:i]
		}
	}
	return []string{}
}
