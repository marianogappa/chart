package main

import (
	"fmt"
	"log"
	"text/template"
	"time"
)

type scatterTemplateData struct {
	FSS       [][]float64
	TSS       [][]time.Time
	SSS       [][]string
	MinFSS    []float64
	MaxFSS    []float64
	Colors    map[string]string
	Title     string
	ScaleType string
	XLabel    string
	YLabel    string
}

var scatterTemplate, scatterWithLabelsTemplate *template.Template

func init() {
	scatterTemplateString := `{
    type: 'bubble',
    data: {
        datasets: [
            {
                label: '',
                data: [
                    {{- $times := .TSS -}}
                    {{- $minFSS := .MinFSS -}}
                    {{- $maxFSS := .MaxFSS -}}
                    {{- if len $times -}}
                        {{- $fss := .FSS -}}
                        {{- range $i,$v := $times}}{{if $i}},{{end -}}
                            {{- if len $v -}}
                            {{- $fs := index $fss $i -}}
                              {
                                x: '{{- (index $v 0).Format "2006-01-02T15:04:05.999999999" -}}',
                                y: {{- if ge (len $fs) 1}}{{index $fs 0 | printf "%g"}}{{else}}0{{end}},
                                r: {{- if ge (len $fs) 2}}{{scatterRadius (index $fs 1) (index $minFSS 1) (index $maxFSS 1) | printf "%g"}}{{else}}4{{end -}}
                              }
                            {{- end -}}
                        {{- end -}}
                    {{- else -}}
                        {{- range $i,$v := .FSS}}{{if $i}},{{end -}}
                            {{- if len $v -}}
                              {
                                x: {{- index $v 0 | printf "%g" -}},
                                y: {{- if ge (len $v) 2}}{{index $v 1 | printf "%g"}}{{else}}0{{end}},
                                r: {{- if ge (len $v) 3}}{{scatterRadius (index $v 2) (index $minFSS 2) (index $maxFSS 2) | printf "%g"}}{{else}}4{{end -}}
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
                    var value = data.datasets[tti.datasetIndex].data[tti.index];
                    if (value.y) {
                        value = value.y
                    }
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

	scatterWithLabelsTemplateString := `{
    type: 'bubble',
    data: {
        datasets: [

                    {{- $times := .TSS -}}
                    {{- $minFSS := .MinFSS -}}
                    {{- $maxFSS := .MaxFSS -}}
                    {{- $colors := .Colors -}}
                    {{- $sss := .SSS -}}
                    {{- if len $times -}}
                        {{- $fss := .FSS -}}
                        {{- range $i,$v := $times}}{{if $i}},{{end -}}
                            {{- if len $v -}}
                            {{- $fs := index $fss $i -}}
            {
                label: '{{(index (index $sss $i) 0)}}',
                data: [
                              {
                                x: '{{- (index $v 0).Format "2006-01-02T15:04:05.999999999" -}}',
                                y: {{- if ge (len $fs) 1}}{{index $fs 0 | printf "%g"}}{{else}}0{{end}},
                                r: {{- if ge (len $fs) 2}}{{scatterRadius (index $fs 1) (index $minFSS 1) (index $maxFSS 1) | printf "%g"}}{{else}}4{{end -}}
                              }
                ],
                backgroundColor: {{index $colors (index (index $sss $i) 0)}}
            }
                            {{- end -}}
                        {{- end -}}
                    {{- else -}}
                        {{- range $i,$v := .FSS}}{{if $i}},{{end -}}
                            {{- if len $v -}}
            {
                label: '{{(index (index $sss $i) 0)}}',
                data: [
                              {
                                x: {{- index $v 0 | printf "%g" -}},
                                y: {{- if ge (len $v) 2}}{{index $v 1 | printf "%g"}}{{else}}0{{end}},
                                r: {{- if ge (len $v) 3}}{{scatterRadius (index $v 2) (index $minFSS 2) (index $maxFSS 2) | printf "%g"}}{{else}}4{{end -}}
                              }
                ],
                backgroundColor: {{index $colors (index (index $sss $i) 0)}}
            }
                            {{- end -}}
                        {{- end -}}
                    {{- end -}}
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
		"colorFirstN":   colorFirstN,
		"scatterRadius": scatterRadius,
	}).Parse(scatterTemplateString)
	if err != nil {
		log.Fatal(err)
	}

	scatterWithLabelsTemplate, err = template.New("scatter").Funcs(template.FuncMap{
		"preprocessLabel": preprocessLabel,
		"scatterRadius":   scatterRadius,
	}).Parse(scatterWithLabelsTemplateString)
	if err != nil {
		log.Fatal(err)
	}
}

func setupScatter(fss [][]float64, sss [][]string, tss [][]time.Time, minFSS []float64, maxFSS []float64, title string, scaleType scaleType, xLabel string, yLabel string) (interface{}, *template.Template, error) {
	if len(fss) == 0 {
		return nil, nil, fmt.Errorf("Couldn't find values to plot.")
	}

	css := map[string]string{}
	colorReset()
	for _, ss := range sss {
		if len(ss) > 0 && css[ss[0]] == "" {
			css[ss[0]] = colorNext()
		}
	}

	templateData := scatterTemplateData{
		FSS:       fss,
		TSS:       tss,
		SSS:       sss,
		MinFSS:    minFSS,
		MaxFSS:    maxFSS,
		Colors:    css,
		Title:     title,
		ScaleType: scaleType.string(),
		XLabel:    xLabel,
		YLabel:    yLabel,
	}

	templ := scatterTemplate
	if sss != nil {
		templ = scatterWithLabelsTemplate
	}

	return templateData, templ, nil
}

func scatterRadius(x, min, max float64) float64 {
	if max-min < 50 {
		return x - min + 4
	}
	return float64(4) + (x-min)/(max-min)*50
}
