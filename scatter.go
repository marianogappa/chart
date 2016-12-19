package main

import (
	"fmt"
	"log"
	"text/template"
	"time"
)

type scatterTemplateData struct {
	Floats    [][]float64
	Times     [][]time.Time
	Title     string
	ScaleType string
	XLabel    string
	YLabel    string
}

var scatterTemplate *template.Template

func init() {
	scatterTemplateString := `{
    type: 'bubble',
    data: {
        datasets: [
            {
                label: '',
                data: [
                    {{- $times := .Times -}}
                    {{- if len $times -}}
                        {{- $fss := .Floats -}}
                        {{- range $i,$v := $times}}{{if $i}},{{end -}}
                            {{- if len $v -}}
                            {{- $fs := index $fss $i -}}
                              {
                                x: '{{- (index $v 0).Format "2006-01-02T15:04:05.999999999" -}}',
                                y: {{- if ge (len $fs) 1}}{{index $fs 0 | printf "%g"}}{{else}}0{{end}},
                                r: {{- if ge (len $fs) 2}}{{index $fs 1 | printf "%g"}}{{else}}4{{end -}}
                              }
                            {{- end -}}
                        {{- end -}}
                    {{- else -}}
                        {{- range $i,$v := .Floats}}{{if $i}},{{end -}}
                            {{- if len $v -}}
                              {
                                x: {{- index $v 0 | printf "%g" -}},
                                y: {{- if ge (len $v) 2}}{{index $v 1 | printf "%g"}}{{else}}0{{end}},
                                r: {{- if ge (len $v) 3}}{{index $v 2 | printf "%g"}}{{else}}4{{end -}}
                              }
                            {{- end -}}
                        {{- end -}}
                    {{- end -}}
                ],
                backgroundColor: {{colorFirstN 1}}
            }
        ]
    },
    options: {
        title: {
            display: {{ if len .Title }}true{{else}}false{{end}},
            text: '{{ .Title }}'
        },
        legend: {
            display: false
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
         scales: {
            yAxes: [{
                type: "{{ .ScaleType }}",
                ticks: {
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
                {{if len $times }}type: 'time',{{end}}
                scaleLabel: {
                    display: {{if eq .XLabel ""}}false{{else}}true{{end}},
                    labelString: '{{ .XLabel }}'
                }
            }]
        }
    }
}`

	var err error
	scatterTemplate, err = template.New("scatter").Funcs(template.FuncMap{
		"preprocessLabel": preprocessLabel,
		"colorFirstN":     colorFirstN,
	}).Parse(scatterTemplateString)
	if err != nil {
		log.Fatal(err)
	}
}

func setupScatter(fss [][]float64, sss [][]string, tss [][]time.Time, title string, scaleType scaleType, xLabel string, yLabel string) (interface{}, *template.Template, error) {
	if len(fss) == 0 {
		return nil, nil, fmt.Errorf("Couldn't find values to plot.")
	}

	templateData := scatterTemplateData{
		Floats:    fss,
		Times:     tss,
		Title:     title,
		ScaleType: scaleType.string(),
		XLabel:    xLabel,
		YLabel:    yLabel,
	}

	return templateData, scatterTemplate, nil
}
