package main

import (
	"fmt"
	"log"
	"text/template"
)

type pieTemplateData struct {
	Labels    [][]string
	Title     string
	ChartType string
	Data      [][]float64
}

var pieTemplate *template.Template

func init() {
	pieTemplateString := `{
    type: '{{ .ChartType }}',
    data: {
        labels: [{{ if len .Labels }}{{ range $i,$v := .Labels }}{{if $i}},{{end}}{{if len $v}}{{index $v 0 | preprocessLabel}}{{else}}''{{end}}{{end}}{{end}}],
        datasets: [{
            data: [{{ range $i,$v := .Data }}{{if $i}},{{end}}{{if len $v}}{{index $v 0 | printf "%g"}}{{else}}0{{end}}{{end}}],
            backgroundColor: [{{ len .Data | colorFirstN }}]
        }]
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
                    var total = data.datasets[0].data.reduce((a, b) => a + b, 0)
                    var label = data.labels[tti.index];
                    var percentage = Math.round(value / total * 100);
                    return {{ if len .Labels }}label + ': ' + percentage + '%'{{else}}percentage + '%'{{end}};
                }
            }
        }
    }
}`

	var err error
	pieTemplate, err = template.New("pie").Funcs(template.FuncMap{
		"preprocessLabel": preprocessLabel,
		"colorFirstN":     colorFirstN,
	}).Parse(pieTemplateString)
	if err != nil {
		log.Fatal(err)
	}
}

func setupPie(fss [][]float64, sss [][]string, title string) (interface{}, *template.Template, error) {
	if len(fss) == 0 || len(sss) == 0 {
		return nil, nil, fmt.Errorf("Couldn't find values to plot.")
	}

	templateData := pieTemplateData{
		ChartType: "pie",
		Data:      fss,
		Labels:    sss,
		Title:     title,
	}

	return templateData, pieTemplate, nil
}
