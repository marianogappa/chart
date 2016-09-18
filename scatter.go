package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/template"
)

type scatterTemplateData struct {
	Datasets        string
	DisplayTitle    bool
	Title           string
	ChartType       string
	TooltipTemplate string
	ScaleType       string
	XLabel          string
	YLabel          string
}

type scatterDatasetTemplateData struct {
	Label string
	Data  string
	Color string
}

type scatterDatapointTemplateData struct {
	X string
	Y string
	R string
}

var scatterTemplate *template.Template
var scatterDatasetTemplate *template.Template
var scatterDatapointTemplate *template.Template

func init() {
	scatterDatapointTemplateString := `
                {
                    x: {{ .X }},
                    y: {{ .Y }},
                    r: {{ .R }}
                }
	`

	scatterDatasetTemplateString := `
            {
                label: '{{ .Label }}',
                data: [{{ .Data }}],
                backgroundColor: {{.Color}}
            }
	`

	scatterTemplateString := `{
    type: '{{ .ChartType }}',
    data: {
        datasets: [
            {{ .Datasets }}
        ]
    },
    options: {
        title: {
            display: {{ .DisplayTitle }},
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
                    return {{ .TooltipTemplate }};
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
	scatterTemplate, err = template.New("scatter").Parse(scatterTemplateString)
	if err != nil {
		log.Fatal(err)
	}

	scatterDatasetTemplate, err = template.New("scatter_dataset").Parse(scatterDatasetTemplateString)
	if err != nil {
		log.Fatal(err)
	}

	scatterDatapointTemplate, err = template.New("scatter_datapoint").Parse(scatterDatapointTemplateString)
	if err != nil {
		log.Fatal(err)
	}
}

func setupScatter(fss [][]float64, sss [][]string, title string, scaleType scaleType, xLabel string, yLabel string) (interface{}, *template.Template, error) {
	if len(fss) == 0 {
		return nil, nil, fmt.Errorf("Couldn't find values to plot.")
	}

	var ds []string
	for _, fs := range fss {
		if len(fs) < 2 {
			log.Printf("Ignoring line %v; less than 2 coordinates.\n", fs)
			continue
		}
		datapoint := scatterDatapointTemplateData{}
		datapoint.X = strconv.FormatFloat(fs[0], 'f', -1, 64)
		datapoint.Y = strconv.FormatFloat(fs[1], 'f', -1, 64)
		if len(fs) >= 3 {
			datapoint.R = strconv.FormatFloat(fs[2], 'f', -1, 64)
		} else {
			datapoint.R = "4"
		}
		var b bytes.Buffer
		scatterDatapointTemplate.Execute(&b, datapoint)
		ds = append(ds, b.String())
	}

	stringData := strings.Join(ds, ",")

	var b bytes.Buffer
	scatterDatasetTemplate.Execute(&b, scatterDatasetTemplateData{Data: stringData, Color: colorFirstN(1)})
	dataset := b.String()

	templateData := scatterTemplateData{
		ChartType:       "bubble",
		Datasets:        dataset,
		Title:           title,
		DisplayTitle:    len(title) > 0,
		TooltipTemplate: `value`,
		ScaleType:       scaleType.string(),
		XLabel:          xLabel,
		YLabel:          yLabel,
	}

	return templateData, scatterTemplate, nil
}
