package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"

	"../chart/color"
)

type pieTemplateData struct {
	Labels       string
	DisplayTitle bool
	Title        string
	ChartType    string
	Data         string
	Colors       string
}

var pieTemplate *template.Template

func init() {
	pieTemplateString := `{
    type: '{{ .ChartType }}',
    data: {
		labels: [{{ .Labels }}],
        datasets: [{
            data: [{{ .Data }}],
            backgroundColor: [{{ .Colors }}]
        }]
    },
    options: {
        title: {
            display: {{ .DisplayTitle }},
            text: '{{ .Title }}'
        }
    }
}`

	var err error
	pieTemplate, err = template.New("pie").Parse(pieTemplateString)
	if err != nil {
		log.Fatal(err)
	}
}

func setupPie(title string, displayTitle bool, separator rune) (string, error) {
	r := csv.NewReader(os.Stdin)
	r.Comma = separator
	r.Comment = '#'

	var labels []string
	var data []float64
	var d string

	for {
		var fS, s string
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		if len(record) == 0 {
			break
		}
		if len(record) == 1 {
			fS = record[0]
			s = ""
		}
		if len(record) >= 2 {
			s = record[0]
			fS = record[1]
		}
		f, err := strconv.ParseFloat(fS, 64)
		if err != nil {
			log.Printf("Ignoring this as it's not a number: [%v]", d)
			continue
		}
		data = append(data, f)
		labels = append(labels, `"`+s+`"`)
	}

	var ds []string
	for _, v := range data {
		ds = append(ds, strconv.FormatFloat(v, 'f', -1, 64))
	}
	stringData := strings.Join(ds, ",")
	stringLabels := strings.Join(labels, ",")

	templateData := pieTemplateData{
		ChartType:    "pie",
		Data:         stringData,
		Labels:       stringLabels,
		Title:        title,
		DisplayTitle: displayTitle,
		Colors:       color.FirstN(len(stringData)),
	}

	var b bytes.Buffer
	if err := pieTemplate.Execute(&b, templateData); err != nil {
		log.Fatal(err)
	}

	return b.String(), nil
}
