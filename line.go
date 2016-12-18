package main

import (
	"fmt"
	"log"
	"text/template"
)

type lineTemplateData struct {
	Labels    [][]string
	Datasets  [][]float64
	Title     string
	ScaleType string
	XLabel    string
	YLabel    string
}

var lineTemplate *template.Template

func init() {
	lineTemplateString := `{
    type: 'line',
    data: {
        labels: [{{ if len .Labels }}{{ range $i,$v := .Labels }}{{if $i}},{{end}}{{if len $v}}{{index $v 0 | preprocessLabel}}{{else}}''{{end}}{{end}}{{end}}],
        datasets: [{{$datasets := .Datasets}}
            {{range $i, $unused := (index $datasets 0)}}{{if $i}},{{end}}
                {
                  fill: false,
                  data: [
                      {{range $j, $v := $datasets}}{{if $j}},{{end}}{{if gt (len $v) $i}}{{index $v $i | printf "%g"}}{{else}}0{{end}}{{end}}
                  ],
                  borderColor: [{{colorIndex $i}}]
                }
            {{end}}
        ]
    },
    options: {
        title: {
            display: {{ if len .Title }}true{{else}}false{{end}},
            text: '{{ .Title }}'
        },
        tooltips: {
            callbacks: {
                label: function(tti, data) {
                    var value = data.datasets[0].data[tti.index];
                    var label = data.labels[tti.index];
                    return value;
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
                },
                scaleLabel: {
                    display: {{if eq .YLabel ""}}false{{else}}true{{end}},
                    labelString: '{{ .YLabel }}'
                }
            }],
            xAxes: [{
                scaleLabel: {
                    display: {{if eq .XLabel ""}}false{{else}}true{{end}},
                    labelString: '{{ .XLabel }}'
                }
            }]
        }
    }
}`

	var err error
	lineTemplate, err = template.New("line").Funcs(template.FuncMap{
		"preprocessLabel": preprocessLabel,
		"colorIndex":      colorIndex,
	}).Parse(lineTemplateString)
	if err != nil {
		log.Fatal(err)
	}
}

func setupLine(fss [][]float64, sss [][]string, title string, scaleType scaleType, xLabel string, yLabel string) (interface{}, *template.Template, error) {
	if len(fss) == 0 || len(sss) == 0 {
		return nil, nil, fmt.Errorf("Couldn't find values to plot.")
	}

	templateData := lineTemplateData{
		Datasets:  fss,
		Labels:    sss,
		Title:     title,
		ScaleType: scaleType.string(),
		XLabel:    xLabel,
		YLabel:    yLabel,
	}

	return templateData, lineTemplate, nil
}
