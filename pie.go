package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/template"
)

type pieTemplateData struct {
	Labels          string
	DisplayTitle    bool
	Title           string
	ChartType       string
	Data            string
	Colors          string
	TooltipTemplate string
}

var pieTemplate *template.Template

func init() {
	pieTemplateString := `{
    type: '{{ .ChartType }}',
    data: {
		labels: [{{ .Labels }}],
        datasets: [{
            data: [{{ .Data }}],
            backgroundColor: [{{ .Colors }}]
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
                    var total = data.datasets[0].data.reduce((a, b) => a + b, 0)
                    var label = data.labels[tti.index];
                    var percentage = Math.round(value / total * 100);
                    return {{ .TooltipTemplate }};
                }
            }
        }
    }
}`

	var err error
	pieTemplate, err = template.New("pie").Parse(pieTemplateString)
	if err != nil {
		log.Fatal(err)
	}
}

func setupPie(fss [][]float64, sss [][]string, title string, displayTitle bool) (string, error) {

	var ds []string
	for _, fs := range fss {
		if len(fs) == 0 {
			return "", fmt.Errorf("Couldn't find values to plot.") //TODO this probably shouldn't happen
		}
		ds = append(ds, strconv.FormatFloat(fs[0], 'f', -1, 64))
	}

	var ls []string

	noLabels := len(sss) == 0
	for _, ss := range sss {
		if len(ss) == 0 {
			noLabels = true //TODO this probably shouldn't happen
			break
		}
		ls = append(ls, "`"+ss[0]+"`")
	}

	stringData := strings.Join(ds, ",")
	stringLabels := strings.Join(ls, ",")

	var tooltipTemplate string
	if noLabels {
		tooltipTemplate = `percentage + '%'`
	} else {
		tooltipTemplate = `label + ': ' + percentage + '%'`
	}

	templateData := pieTemplateData{
		ChartType:       "pie",
		Data:            stringData,
		Labels:          stringLabels,
		Title:           title,
		DisplayTitle:    displayTitle,
		Colors:          colorFirstN(len(stringData)),
		TooltipTemplate: tooltipTemplate,
	}

	var b bytes.Buffer
	if err := pieTemplate.Execute(&b, templateData); err != nil {
		return "", err
	}

	return b.String(), nil
}
