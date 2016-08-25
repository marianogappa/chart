package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"strings"
	"text/template"
)

type pieTemplateData struct {
	Labels          string
	DisplayTitle    bool
	Title           string
	ChartType       string
	Data            string
	Colors          string
	TooltipTemplate string
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
        },
        tooltips: {
            callbacks: {
                label: function(tti, data) {
                    var value = data.datasets[0].data[tti.index];
                    var total = data.datasets[0].data.reduce((a, b) => a + b, 0)
                    var label = data.labels[tti.index];
                    var percentage = Math.round(value / total * 100);
                    return {{ .TooltipTemplate }};
                }
            }
        }
    }
}`

	var err error
	pieTemplate, err = template.New("pie").Parse(pieTemplateString)
	if err != nil {
		log.Fatal(err)
	}
}

func setupPie(input []string, title string, displayTitle bool, separator rune, invert bool) (string, error) {
	i := strings.Join(input, "\n")

	r := csv.NewReader(strings.NewReader(i))
	r.Comma = separator
	r.Comment = '#'

	var labels []string
	var data []float64
	var d, tooltipTemplate string
	var noLabels bool

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
			noLabels = true
		}
		if len(record) >= 2 {
			if invert {
				s = record[1]
				fS = record[0]
			} else {
				s = record[0]
				fS = record[1]
			}
		}
		f, err := strconv.ParseFloat(fS, 64)
		if err != nil {
			log.Printf("Ignoring this as it's not a number: [%v]", d)
			continue
		}
		data = append(data, f)
		if !noLabels {
			labels = append(labels, `"`+s+`"`)
		}
	}

	var ds []string
	for _, v := range data {
		ds = append(ds, strconv.FormatFloat(v, 'f', -1, 64))
	}
	stringData := strings.Join(ds, ",")
	stringLabels := strings.Join(labels, ",")

	if noLabels {
		tooltipTemplate = `percentage + '%'`
	} else {
		tooltipTemplate = `label + ': ' + percentage + '%'`
	}

	templateData := pieTemplateData{
		ChartType:       "pie",
		Data:            stringData,
		Labels:          stringLabels,
		Title:           title,
		DisplayTitle:    displayTitle,
		Colors:          colorFirstN(len(stringData)),
		TooltipTemplate: tooltipTemplate,
	}

	var b bytes.Buffer
	if err := pieTemplate.Execute(&b, templateData); err != nil {
		log.Fatal(err)
	}

	return b.String(), nil
}
