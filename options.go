package main

import (
	"flag"
	"log"
)

type chartType int
type scaleType int

const (
	pie chartType = iota
	bar
)
const (
	linear scaleType = iota
	logarithmic
)

type options struct {
	title     string
	separator string
	scaleType scaleType
	chartType chartType
}

var defaultOptions = options{
	title:     "",
	separator: "\t",
	scaleType: linear,
	chartType: pie,
}

func mustResolveOptions(args []string) options {
	o, err := resolveOptions(args)
	if err != nil {
		log.Fatal(err)
	}
	return o
}

func resolveOptions(args []string) (options, error) {
	titleHelp := "Sets the title for the chart."
	separatorHelp := "Sets the separator for each row's fields; can be ' ', '\\t', ';', ','; default \\t."
	logarithmicHelp := "Sets logarithmic scale for the y-axis."
	chartTypeHelp := "Sets the chart type; default is 'pie'; can be pie/bar."

	o := defaultOptions

	var chartType string
	var log bool

	fs := flag.NewFlagSet("params", flag.ContinueOnError)
	fs.StringVar(&o.title, "title", o.title, titleHelp)
	fs.StringVar(&o.title, "t", o.title, titleHelp)
	fs.StringVar(&o.separator, "separator", o.separator, separatorHelp)
	fs.StringVar(&o.separator, "s", o.separator, separatorHelp)
	fs.BoolVar(&log, "log", o.scaleType.isLogarithmic(), logarithmicHelp)
	fs.BoolVar(&log, "l", o.scaleType.isLogarithmic(), logarithmicHelp)
	fs.StringVar(&chartType, "type", o.chartType.string(), chartTypeHelp)
	fs.StringVar(&chartType, "y", o.chartType.string(), chartTypeHelp)

	err := fs.Parse(fromFirstDash(args))
	if err != nil {
		return o, err
	}

	o.chartType.become(chartType)
	o.scaleType.become(log)

	if o.separator != " " && o.separator != ";" && o.separator != "," && o.separator != "\t" {
		o.separator = defaultOptions.separator
	}

	hints := fs.Args()
	hints = append(hints, untilFirstDash(args)...)
	for _, h := range hints {
		switch h {
		case "log":
			o.scaleType = logarithmic
		case "bar":
			o.chartType = bar
		case ",":
			o.separator = ","
		case ";":
			o.separator = ";"
		case " ":
			o.separator = " "
		}
	}

	return o, nil
}

func (c chartType) string() string {
	if c == bar {
		return "bar"
	}
	return "pie"
}

func (c *chartType) become(s string) {
	if s == "bar" {
		*c = bar
	}
}

func (c scaleType) isLogarithmic() bool {
	return c == logarithmic
}

func (c *scaleType) become(log bool) {
	if log {
		*c = logarithmic
	}
}

func (c scaleType) string() string {
	if c == logarithmic {
		return "logarithmic"
	}
	return "linear"
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
	return as
}
