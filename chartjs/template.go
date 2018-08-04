package chartjs

import (
	"log"
	"text/template"
)

var cjsTemplate *template.Template

func init() {
	cjsTemplateString :=
		`{
    type: '{{ .ChartType }}',{{ $manyColor := or (eq .ChartType "pie") (eq .ChartType "bar") }}
    data: {
      labels: [{{ .Labels }}],
      datasets: [
        {{range $i,$v := .Datasets}}{{if $i}},{{end -}}
        {
          fill: {{ .Fill }},
          {{if len .Label}}label: '{{ .Label }}',{{end}}
          {{if len .BackgroundColor}}backgroundColor: {{if $manyColor}}[{{end}}{{ .BackgroundColor }}{{if $manyColor}}]{{end}},{{end}}
          {{if len .BorderColor}}borderColor: {{ .BorderColor }},{{end}}
          data: [
            {{if len .SimpleData}}{{range $i,$v := .SimpleData}}{{if $i}},{{end -}}{{.}}{{end}}{{end}}
            {{if len .ComplexData}}{{range $i,$v := .ComplexData}}{{if $i}},{{end -}}
              {
                x: {{ .X }},
                y: {{ .Y -}}
                {{- if .UsesR}},
                r: {{ .R }}
                {{end}}
              }
            {{end}}{{end}}
          ]
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
                {{ .TooltipCallback }}
              }
          }
      },
      {{ if ne .ChartType "pie" }}
        legend: {
            display: {{ if gt (len .Datasets) 1 }}true,
            position: 'bottom'{{else}}false{{end}}
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
                {{ if .UsesTimeScale }}
                type: 'time',
                position: 'bottom',
                {{ else if eq .ActualChartType "scatterline" }}
                type: 'linear',
                position: 'bottom',
                {{end}}
                scaleLabel: {
                    display: {{if eq .XLabel ""}}false{{else}}true{{end}},
                    labelString: '{{ .XLabel }}'
                }
            }]
        },
        elements: {
            line: {
                tension: 0, // disables bezier curves
            },
        },
        {{end}}
        animation: {
            duration: 0, // general animation time
        },
        hover: {
            animationDuration: 0, // duration of animations when hovering an item
        },
        responsiveAnimationDuration: 0, // animation duration after a resize
    }
}`

	var err error
	cjsTemplate, err = template.New("").Parse(cjsTemplateString)
	if err != nil {
		log.Fatal(err)
	}
}
