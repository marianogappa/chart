package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/template"
)

type lineTemplateData struct {
	Labels          string
	Datasets        string
	DisplayTitle    bool
	Title           string
	ChartType       string
	TooltipTemplate string
	ScaleType       string
}

type lineDatasetTemplateData struct {
	Data  string
	Color string
}

var lineTemplate *template.Template
var lineDatasetTemplate *template.Template

func init() {
	lineDatasetTemplateString := `
            {
                fill: false,
                data: [{{ .Data }}],
                borderColor: [{{.Color}}]
            }
	`

	lineTemplateString := `{
    type: '{{ .ChartType }}',
    data: {
		labels: [{{ .Labels }}],
        datasets: [
            {{ .Datasets }}
        ]
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
	lineTemplate, err = template.New("line").Parse(lineTemplateString)
	if err != nil {
		log.Fatal(err)
	}

	lineDatasetTemplate, err = template.New("line_dataset").Parse(lineDatasetTemplateString)
	if err != nil {
		log.Fatal(err)
	}
}

func setupLine(fss [][]float64, sss [][]string, title string, scaleType scaleType) (interface{}, *template.Template, error) {
	if len(fss) == 0 || len(sss) == 0 {
		return nil, nil, fmt.Errorf("Couldn't find values to plot.")
	}

	var ds [][]string
	for _, fs := range fss {
		if len(fs) == 0 {
			return nil, nil, fmt.Errorf("Couldn't find values to plot.") //TODO this probably shouldn't happen
		}
		for i, f := range fs {
			if i > len(ds)-1 {
				ds = append(ds, []string{})
			}
			ds[i] = append(ds[i], strconv.FormatFloat(f, 'f', -1, 64))
		}
		ds = append(ds)
	}

	var ls []string
	for _, ss := range sss {
		if len(ss) == 0 {
			break //TODO this probably shouldn't happen
		}
		ls = append(ls, "`"+ss[0]+"`")
	}

	datasets := []string{}
	for i, d := range ds {
		var b bytes.Buffer
		if err := lineDatasetTemplate.Execute(&b, lineDatasetTemplateData{
			Data:  strings.Join(d, ","),
			Color: colorIndex(i),
		}); err != nil {
			return nil, nil, fmt.Errorf("Could not prepare ChartJS js code for chart: [%v]", err)
		}

		datasets = append(datasets, b.String())
	}

	stringDatasets := strings.Join(datasets, ",")
	stringLabels := strings.Join(ls, ",")

	templateData := lineTemplateData{
		ChartType:       "line",
		Datasets:        stringDatasets,
		Labels:          stringLabels,
		Title:           title,
		DisplayTitle:    len(title) > 0,
		TooltipTemplate: `value`,
		ScaleType:       scaleType.string(),
	}

	return templateData, lineTemplate, nil
}
