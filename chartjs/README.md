# chart/chartjs [![Build Status](https://img.shields.io/travis/marianogappa/chart.svg)](https://travis-ci.org/marianogappa/chart) [![Coverage Status](https://coveralls.io/repos/github/MarianoGappa/chart/badge.svg?branch=master)](https://coveralls.io/github/MarianoGappa/chart?branch=master) [![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/marianogappa/chart/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/marianogappa/chart?style=flat-square)](https://goreportcard.com/report/github.com/marianogappa/chart) [![GoDoc](https://godoc.org/github.com/marianogappa/chart/chartjs?status.svg)](https://godoc.org/github.com/marianogappa/chart/chartjs)

## Example usage

```go
package main

import (
	"log"
	"os"
	"strings"

	"github.com/marianogappa/chart/chartjs"
	"github.com/marianogappa/chart/dataset"
	"github.com/marianogappa/chart/format"
)

func main() {
	fileContent := `Category1	1
Category2	2
Category3	3
`
	reader := strings.NewReader(fileContent)
	lineFormat, _ := format.NewLineFormat("sf", '\t', "")        // Look into format.Parse to infer the line format
	ds, err := dataset.New(reader, lineFormat)                   // Construct dataset manually if not reading a file
	if err != nil {
		log.Fatal(err)
	}
	chart := chartjs.New(chartjs.Pie, *ds, chartjs.Options{Title: "Example chart"}) // Consult godoc for ChartTypes
	if err := chart.Build(chartjs.OutputAll, os.Stdout); err != nil {               // Consult godoc for OutputModes
		log.Fatal(err)
	}
}
```

Try it like this

```
$ go run main.go > /tmp/a.html && open /tmp/a.html
```
