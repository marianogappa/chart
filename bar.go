package main

import (
	"bytes"
	"fmt"
	"log"
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

func setupBar(fss [][]float64, sss [][]string, title string, scaleType scaleType) (string, error) {

	var ds []string
	for _, fs := range fss {
		if len(fs) == 0 {
			return "", fmt.Errorf("Couldn't find values to plot.") //TODO this probably shouldn't happen
		}
		ds = append(ds, strconv.FormatFloat(fs[0], 'f', -1, 64))
	}

	var ls []string
	for _, ss := range sss {
		if len(ss) == 0 {
			break //TODO this probably shouldn't happen
		}
		ls = append(ls, "`"+ss[0]+"`")
	}

	stringData := strings.Join(ds, ",")
	stringLabels := strings.Join(ls, ",")

	templateData := barTemplateData{
		ChartType:       "bar",
		Data:            stringData,
		Labels:          stringLabels,
		Title:           title,
		DisplayTitle:    len(title) > 0,
		Colors:          colorFirstN(len(stringData)),
		TooltipTemplate: `value`,
		ScaleType:       scaleType.string(),
	}

	var b bytes.Buffer
	if err := barTemplate.Execute(&b, templateData); err != nil {
		return "", err
	}

	return b.String(), nil
}
