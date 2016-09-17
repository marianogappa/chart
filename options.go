package main

import (
	"flag"
	"log"
)

type chartType int
type scaleType int

const (
	undefinedChartType chartType = iota
	pie
	bar
	line
)
const (
	linear scaleType = iota
	logarithmic
)

type options struct {
	title     string
	separator rune
	scaleType scaleType
	chartType chartType
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

	o := options{}

	fs := flag.NewFlagSet("params", flag.ContinueOnError)
	fs.StringVar(&o.title, "title", o.title, titleHelp)
	fs.StringVar(&o.title, "t", o.title, titleHelp)

	err := fs.Parse(fromFirstDash(args))
	if err != nil {
		return o, err
	}

	hints := fs.Args()
	hints = append(hints, untilFirstDash(args)...)
	for _, h := range hints {
		switch h {
		case "log":
			o.scaleType = logarithmic
		case "bar":
			o.chartType = bar
		case "pie":
			o.chartType = pie
		case "line":
			o.chartType = line
		case ",":
			o.separator = ','
		case ";":
			o.separator = ';'
		case " ":
			o.separator = ' '
		}
	}

	if o.separator != ' ' && o.separator != ';' && o.separator != ',' {
		o.separator = '\t'
	}

	return o, nil
}

func (c chartType) string() string {
	switch c {
	case bar:
		return "bar"
	case line:
		return "line"
	default:
		return "pie"
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
