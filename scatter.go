package main

import (
	"fmt"
	"log"
	"text/template"
)

type scatterTemplateData struct {
	Data      [][]float64
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
                data: [{{range $i,$v := .Data}}{{if $i}},{{end -}}
                {{- if len $v -}}
                  {
                    x: {{- index $v 0 | printf "%g" -}},
                    y: {{- if ge (len $v) 2}}{{index $v 1 | printf "%g"}}{{else}}0{{end}},
                    r: {{- if ge (len $v) 3}}{{index $v 2 | printf "%g"}}{{else}}4{{end -}}
                  }
                {{- end -}}
                {{- end}}],
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

func setupScatter(fss [][]float64, sss [][]string, title string, scaleType scaleType, xLabel string, yLabel string) (interface{}, *template.Template, error) {
	if len(fss) == 0 {
		return nil, nil, fmt.Errorf("Couldn't find values to plot.")
	}

	templateData := scatterTemplateData{
		Data:      fss,
		Title:     title,
		ScaleType: scaleType.string(),
		XLabel:    xLabel,
		YLabel:    yLabel,
	}

	return templateData, scatterTemplate, nil
}
