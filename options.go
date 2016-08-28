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

	o := defaultOptions

	fs := flag.NewFlagSet("params", flag.ContinueOnError)
	fs.StringVar(&o.title, "title", o.title, titleHelp)
	fs.StringVar(&o.title, "t", o.title, titleHelp)

	err := fs.Parse(fromFirstDash(args))
	if err != nil {
		return o, err
	}

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
