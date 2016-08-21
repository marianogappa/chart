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
)

type barTemplateData struct {
	Labels          string
	DisplayTitle    bool
	Title           string
	ChartType       string
	Data            string
	Colors          string
	TooltipTemplate string
	ScaleType       string
}

var barTemplate *template.Template

func init() {
	barTemplateString := `{
    type: '{{ .ChartType }}',
    data: {
		labels: [{{ .Labels }}],
        datasets: [{
            label: "My First dataset",
            data: [{{ .Data }}],
            backgroundColor: [{{.Colors}}]
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
                    var label = data.labels[tti.index];
                    return {{ .TooltipTemplate }};
                }
            }
        },
        legend: {
            display: false
        },
        scales: {
            yAxes: [{
                type: "{{ .ScaleType }}",
                ticks: {
                    beginAtZero: true,
                    callback: function(value, index, values) {
                        return value;
                    }
                }
            }]
        }
    }
}`

	var err error
	barTemplate, err = template.New("bar").Parse(barTemplateString)
	if err != nil {
		log.Fatal(err)
	}
}

func setupBar(title string, displayTitle bool, separator rune, scaleType string) (string, error) {
	r := csv.NewReader(os.Stdin)
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
			s = record[0]
			fS = record[1]
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

	tooltipTemplate = `value`

	templateData := barTemplateData{
		ChartType:       "bar",
		Data:            stringData,
		Labels:          stringLabels,
		Title:           title,
		DisplayTitle:    displayTitle,
		Colors:          colorFirstN(len(stringData)),
		TooltipTemplate: tooltipTemplate,
		ScaleType:       scaleType,
	}

	var b bytes.Buffer
	if err := barTemplate.Execute(&b, templateData); err != nil {
		log.Fatal(err)
	}

	return b.String(), nil
}
