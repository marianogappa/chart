package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/marianogappa/chart/format"
)

type chartType int
type scaleType int
type colorType int

const (
	undefinedChartType chartType = iota
	pie
	bar
	line
	scatter
)
const (
	linear scaleType = iota
	logarithmic
)
const (
	defaultColor colorType = iota
	legacyColor
	gradient
)

type options struct {
	title         string
	separator     rune
	scaleType     scaleType
	chartType     chartType
	colorType     colorType
	xLabel        string
	yLabel        string
	dateFormat    string
	rawLineFormat string
	lineFormat    format.LineFormat
	zeroBased     bool
	debug         bool
	help          bool
}

func mustResolveOptions(args []string) options {
	o, err := resolveOptions(args)
	if err != nil {
		log.Fatal(err)
	}
	return o
}

func resolveOptions(args []string) (options, error) {
	var (
		titleHelp      = "Sets the title for the chart."
		xLabelHelp     = "Sets the label for the x axis."
		yLabelHelp     = "Sets the label for the y axis."
		dateFormatHelp = "Sets the date format, according to https://golang.org/src/time/format.go."
		lineFormatHelp = "Prevents chart from inferring line format. Syntax: `[dfs]*` where d=date, f=float, s=string. e.g. `df`"
		zeroBasedHelp  = "Makes y-axis begin at zero."
		debugHelp      = "Use to make sure to double-check the chart is showing what you expect."
		helpHelp       = "Show help."
	)

	o := options{}

	fs := flag.NewFlagSet("params", flag.ContinueOnError)
	fs.StringVar(&o.title, "title", o.title, titleHelp)
	fs.StringVar(&o.title, "t", o.title, titleHelp)
	fs.StringVar(&o.xLabel, "x", o.xLabel, xLabelHelp)
	fs.StringVar(&o.yLabel, "y", o.yLabel, yLabelHelp)
	fs.StringVar(&o.dateFormat, "date-format", o.dateFormat, dateFormatHelp)      //TODO document
	fs.StringVar(&o.rawLineFormat, "lineFormat", o.rawLineFormat, lineFormatHelp) //TODO document
	fs.BoolVar(&o.zeroBased, "zero-based", o.zeroBased, zeroBasedHelp)            //TODO document
	fs.BoolVar(&o.debug, "debug", o.debug, debugHelp)                             //TODO document
	fs.BoolVar(&o.help, "help", o.help, helpHelp)
	fs.BoolVar(&o.help, "h", o.help, helpHelp)

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
		case "scatter":
			o.chartType = scatter
		case "legacy-color":
			o.colorType = legacyColor
		case "gradient":
			o.colorType = gradient
		case ",":
			o.separator = ','
		case ";":
			o.separator = ';'
		case " ":
			o.separator = ' '
		case "ANSIC": //TODO document
			o.dateFormat = "Mon Jan _2 15:04:05 2006"
		case "UnixDate":
			o.dateFormat = "Mon Jan _2 15:04:05 MST 2006"
		case "RubyDate":
			o.dateFormat = "Mon Jan 02 15:04:05 -0700 2006"
		case "RFC822":
			o.dateFormat = "02 Jan 06 15:04 MST"
		case "RFC822Z":
			o.dateFormat = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
		case "RFC850":
			o.dateFormat = "Monday, 02-Jan-06 15:04:05 MST"
		case "RFC1123":
			o.dateFormat = "Mon, 02 Jan 2006 15:04:05 MST"
		case "RFC1123Z":
			o.dateFormat = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
		case "RFC3339":
			o.dateFormat = "2006-01-02T15:04:05Z07:00"
		case "RFC3339Nano":
			o.dateFormat = "2006-01-02T15:04:05.999999999Z07:00"
		case "Kitchen":
			o.dateFormat = "3:04PM"
		case "Stamp":
			o.dateFormat = "Jan _2 15:04:05"
		case "StampMilli":
			o.dateFormat = "Jan _2 15:04:05.000"
		case "StampMicro":
			o.dateFormat = "Jan _2 15:04:05.000000"
		case "mysql":
			o.dateFormat = "2006-01-02 15:04:05"
		case "debug":
			o.debug = true
		case "help":
			o.help = true
		}
	}

	if o.help {
		fs.PrintDefaults()
		fmt.Println(`
  pie
	render a pie chart
  bar
	render a bar chart
  line
	render a line chart
  scatter
	render a scatter plot chart
  legacy-color
    use legacy colors
  gradient
  	use color gradients
  log
	use logarithmic scale (bar chart only)
  ' '|';'|','|'\t'
	this character separates columns on each line (\t = default)

More info: https://github.com/marianogappa/chart`)
		return o, fmt.Errorf("Showing usage")
	}

	if o.separator != ' ' && o.separator != ';' && o.separator != ',' {
		o.separator = '\t'
	}

	o.lineFormat, err = format.NewLineFormat(o.rawLineFormat, o.separator, o.dateFormat)
	if err != nil {
		return o, err
	}

	return o, nil
}

func (c chartType) string() string {
	switch c {
	case bar:
		return "bar"
	case line:
		return "line"
	case scatter:
		return "scatter"
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
