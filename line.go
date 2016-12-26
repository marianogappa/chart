package main

import (
	"fmt"
	"log"
	"text/template"
	"time"
)

type lineTemplateData struct {
	SSS       [][]string
	FSS       [][]float64
	TSS       [][]time.Time
	Title     string
	ScaleType string
	XLabel    string
	YLabel    string
}

var lineTemplate, scatterLineTemplate, scatterLineWithTimeTemplate *template.Template

func init() {
	lineTemplateString := `{
    type: 'line',
    data: {
        labels: [{{ if len .SSS }}{{ range $i,$v := .SSS }}{{if $i}},{{end}}{{if len $v}}{{index $v 0 | preprocessLabel}}{{else}}'row {{$i}}'{{end}}{{end}}{{end}}],
        datasets: [{{$fss := .FSS}}
            {{range $i, $unused := (index $fss 0)}}{{if $i}},{{end}}
                {
                  fill: false,
                  data: [
                      {{range $j, $v := $fss}}{{if $j}},{{end}}{{if gt (len $v) $i}}{{index $v $i | printf "%g"}}{{else}}0{{end}}{{end}}
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

	scatterLineTemplateString := `{
    type: 'line',
    data: {
        datasets: [{{$fss := .FSS}}
            {{range $i, $unused := (index $fss 0)}}{{if gt $i 1}},{{end}}{{if $i}}
                {
                  fill: false,
                  label: 'column {{$i}}',
                  data: [
                      {{range $j, $v := $fss}}{{if $j}},{{end}}
                            {
                                x: {{index $v 0 | printf "%g"}},
                                y: {{if gt (len $v) $i}}{{index $v $i | printf "%g"}}{{else}}0{{end}}
                            }
                      {{end}}
                  ],
                  borderColor: [{{colorIndex $i}}]
                }
            {{end}}{{end}}
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
                type: 'linear',
                position: 'bottom',
                scaleLabel: {
                    display: {{if eq .XLabel ""}}false{{else}}true{{end}},
                    labelString: '{{ .XLabel }}'
                }
            }]
        }
    }
}`

	scatterLineWithTimeTemplateString := `{
    type: 'line',
    data: {
        datasets: [{{$fss := .FSS}}{{$tss := .TSS}}
            {{range $i, $unused := (index $fss 0)}}{{if $i}},{{end}}
                {
                  fill: false,
                  label: 'column {{$i}}',
                  data: [
                      {{range $j, $v := $fss}}{{if $j}},{{end}}{{ $t := index $tss $j }}
                            {
                                x: '{{ (index $t 0).Format "2006-01-02T15:04:05.999999999" }}',
                                y: {{if gt (len $v) $i}}{{index $v $i | printf "%g"}}{{else}}0{{end}}
                            }
                      {{end}}
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
                type: 'time',
                position: 'bottom',
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

	scatterLineTemplate, err = template.New("line").Funcs(template.FuncMap{
		"preprocessLabel": preprocessLabel,
		"colorIndex":      colorIndex,
	}).Parse(scatterLineTemplateString)
	if err != nil {
		log.Fatal(err)
	}

	scatterLineWithTimeTemplate, err = template.New("line").Funcs(template.FuncMap{
		"preprocessLabel": preprocessLabel,
		"colorIndex":      colorIndex,
	}).Parse(scatterLineWithTimeTemplateString)
	if err != nil {
		log.Fatal(err)
	}
}

func setupLine(fss [][]float64, sss [][]string, tss [][]time.Time, title string, scaleType scaleType, xLabel string, yLabel string) (interface{}, *template.Template, error) {
	if fss == nil || (sss == nil && tss == nil && len(fss[0]) < 2) {
		return nil, nil, fmt.Errorf("Couldn't find values to plot.")
	}

	templateData := lineTemplateData{
		FSS:       fss,
		SSS:       sss,
		TSS:       tss,
		Title:     title,
		ScaleType: scaleType.string(),
		XLabel:    xLabel,
		YLabel:    yLabel,
	}

	templ := lineTemplate
	if sss == nil && tss == nil {
		templ = scatterLineTemplate
	} else if sss == nil {
		templ = scatterLineWithTimeTemplate
	}

	return templateData, templ, nil
}
