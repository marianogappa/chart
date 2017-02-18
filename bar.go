package main

import (
	"fmt"
	"log"
	"text/template"
	"time"
)

type barTemplateData struct {
	Labels          [][]string
	Title           string
	ChartType       string
	Data            [][]float64
	TooltipTemplate string
	ScaleType       string
	XLabel          string
	YLabel          string
	ZeroBased       bool
}

var barTemplate *template.Template

func init() {
	barTemplateString := `{
    type: '{{ .ChartType }}',
    data: {
        labels: [{{ if len .Labels }}{{ range $i,$v := .Labels }}{{if $i}},{{end}}{{if len $v}}{{index $v 0 | preprocessLabel}}{{else}}'row {{$i}}'{{end}}{{end}}{{end}}],
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
                    beginAtZero: {{ .ZeroBased }},
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
	barTemplate, err = template.New("bar").Funcs(template.FuncMap{
		"preprocessLabel": preprocessLabel,
		"colorFirstN":     colorFirstN,
	}).Parse(barTemplateString)
	if err != nil {
		log.Fatal(err)
	}
}

func setupBar(fss [][]float64, sss [][]string, tss [][]time.Time, title string, scaleType scaleType, xLabel string, yLabel string, zeroBased bool) (interface{}, *template.Template, error) {
	if len(fss) == 0 || (len(fss[0]) == 1 && len(sss) == 0 && len(tss) == 0) {
		return nil, nil, fmt.Errorf("Couldn't find values to plot.")
	}

	if len(sss) == 0 && len(tss) > 0 {
		for _, ts := range tss {
			ss := make([]string, len(ts))
			for i, t := range ts {
				ss[i] = t.Format("2006-01-02T15:04:05.999999999")
			}
			sss = append(sss, ss)
		}
	}

	if len(sss) == 0 && len(tss) == 0 && len(fss[0]) > 1 {
		for i, fs := range fss {
			sss = append(sss, []string{fmt.Sprintf("%g", fs[0])})
			fss[i] = fss[i][1:]
		}
	}

	templateData := barTemplateData{
		ChartType:       "bar",
		Data:            fss,
		Labels:          sss,
		Title:           title,
		TooltipTemplate: `value`,
		ScaleType:       scaleType.string(),
		XLabel:          xLabel,
		YLabel:          yLabel,
		ZeroBased:       zeroBased,
	}

	return templateData, barTemplate, nil
}
