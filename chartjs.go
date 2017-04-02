package main

import (
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"
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
        }
        {{ if ne .ChartType "pie" }},
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
                  {{ if .UsesTimeScale }}
                  type: 'time',
                  position: 'bottom',
                  {{end}}
                  scaleLabel: {
                      display: {{if eq .XLabel ""}}false{{else}}true{{end}},
                      labelString: '{{ .XLabel }}'
                  }
              }]
          }
        {{end}}
    }
}`

	var err error
	cjsTemplate, err = template.New("").Parse(cjsTemplateString)
	if err != nil {
		log.Fatal(err)
	}
}

type inData struct {
	ChartType string
	FSS       [][]float64
	TSS       [][]time.Time
	SSS       [][]string
	MinFSS    []float64
	MaxFSS    []float64
	Title     string
	ScaleType string
	XLabel    string
	YLabel    string
	ZeroBased bool
}

func (i inData) hasFloats() bool  { return len(i.FSS) > 0 }
func (i inData) hasStrings() bool { return len(i.SSS) > 0 }
func (i inData) hasTimes() bool   { return len(i.TSS) > 0 }

type cjsChart struct {
	inData inData
}

func (c cjsChart) chart() (interface{}, *template.Template, error) {
	return c.data(), cjsTemplate, nil
}

type cjsData struct {
	ChartType       string
	Title           string
	ScaleType       string
	XLabel          string
	YLabel          string
	ZeroBased       bool
	Labels          string // Need backticks; can't use array
	Datasets        []cjsDataset
	TooltipCallback string
	UsesTimeScale   bool
}

type cjsDataset struct {
	SimpleData      []string
	ComplexData     []cjsDataPoint
	BackgroundColor string
	Fill            bool
	Label           string
	BorderColor     string
}

type cjsDataPoint struct {
	X, Y, R string
	UsesR   bool
}

func (c cjsChart) data() cjsData {
	d := c.labelsAndDatasets()
	d.ChartType = c.inData.ChartType
	d.Title = c.inData.Title
	d.ScaleType = c.inData.ScaleType
	d.XLabel = c.inData.XLabel
	d.YLabel = c.inData.YLabel
	d.ZeroBased = c.inData.ZeroBased
	d.TooltipCallback = c.tooltipCallback()
	switch d.ChartType {
	case "scatterline":
		d.ChartType = "line"
	case "scatter":
		d.ChartType = "bubble"
	}

	return d
}

func (c cjsChart) labelsAndDatasets() cjsData {
	var usesTimeScale bool
	if c.inData.ChartType == "line" && !c.inData.hasStrings() {
		c.inData.ChartType = "scatterline"
	}
	switch c.inData.ChartType {
	case "pie", "bar":
		return cjsData{
			Labels: c.marshalLabels(),
			Datasets: []cjsDataset{{
				Fill:            true,
				SimpleData:      c.marshalSimpleData(0),
				BackgroundColor: colorFirstN(len(c.inData.FSS)),
			}},
		}
	case "line":
		ds := []cjsDataset{}
		for i := range c.inData.FSS[0] {
			ds = append(ds, cjsDataset{
				Fill:        false,
				SimpleData:  c.marshalSimpleData(i),
				BorderColor: colorIndex(i),
			})
		}
		return cjsData{
			Labels:   c.marshalLabels(),
			Datasets: ds,
		}
	case "scatterline":
		dss := []cjsDataset{}
		for n := range c.inData.FSS[0] {
			ds := []cjsDataPoint{}
			for i := range c.inData.FSS {
				d := cjsDataPoint{}
				if c.inData.hasTimes() {
					usesTimeScale = true
					d.X = "'" + c.inData.TSS[i][0].Format("2006-01-02T15:04:05.999999999") + "'"
					d.Y = fmt.Sprintf("%g", c.inData.FSS[i][n])
				} else {
					d.X = fmt.Sprintf("%g", c.inData.FSS[i][0])
					d.Y = fmt.Sprintf("%g", c.inData.FSS[i][n+1])
				}
				ds = append(ds, d)
			}
			dss = append(dss, cjsDataset{
				Fill:        false,
				Label:       fmt.Sprintf("column %v", n),
				ComplexData: ds,
				BorderColor: colorIndex(n),
			})
		}
		return cjsData{
			Datasets:      dss,
			UsesTimeScale: usesTimeScale,
		}
	case "scatter":
		css := map[string]string{}
		colorReset()
		for _, ss := range c.inData.SSS {
			if len(ss) > 0 && css[ss[0]] == "" {
				css[ss[0]] = colorNext()
			}
		}

		dss := []cjsDataset{}
		for i := range c.inData.FSS {
			d := cjsDataPoint{UsesR: true}
			if c.inData.hasTimes() {
				usesTimeScale = true
				d.X = "'" + c.inData.TSS[i][0].Format("2006-01-02T15:04:05.999999999") + "'"
				d.Y = fmt.Sprintf("%g", c.inData.FSS[i][0])
				if len(c.inData.FSS[i]) >= 2 {
					d.R = fmt.Sprintf("%v", scatterRadius(c.inData.FSS[i][1], c.inData.MinFSS[1], c.inData.MaxFSS[1]))
				} else {
					d.R = fmt.Sprintf("%v", 4)
				}
			} else {
				d.X = fmt.Sprintf("%g", c.inData.FSS[i][0])
				d.Y = fmt.Sprintf("%g", c.inData.FSS[i][1])
				if len(c.inData.FSS[i]) >= 3 {
					d.R = fmt.Sprintf("%v", scatterRadius(c.inData.FSS[i][2], c.inData.MinFSS[2], c.inData.MaxFSS[2]))
				} else {
					d.R = fmt.Sprintf("%v", 4)
				}
			}
			color := colorFirstN(1)
			label := ""
			if c.inData.hasStrings() {
				color = css[c.inData.SSS[i][0]]
				label = c.inData.SSS[i][0]
			}
			dss = append(dss, cjsDataset{
				Fill:            true,
				Label:           label,
				ComplexData:     []cjsDataPoint{d},
				BackgroundColor: color,
			})
		}
		return cjsData{
			Datasets:      dss,
			UsesTimeScale: usesTimeScale,
		}
	default:
		log.Fatalf("Unknown chart type: %v", c.inData.ChartType)
		return cjsData{}
	}
}

func (c cjsChart) marshalLabels() string {
	if !c.inData.hasStrings() && c.inData.hasTimes() {
		ls := make([]string, len(c.inData.TSS))
		for i, ts := range c.inData.TSS {
			ls[i] = ts[0].Format("2006-01-02T15:04:05.999999999")
		}
		return "`" + strings.Join(ls, "`,`") + "`"
	}

	if !c.inData.hasStrings() {
		ls := make([]string, len(c.inData.FSS))
		for i := range c.inData.FSS {
			ls[i] = fmt.Sprintf("slice %v", i)
		}
		return strings.Join(ls, ",")
	}

	ls := make([]string, len(c.inData.SSS))
	for i, l := range c.inData.SSS {
		ls[i] = preprocessLabel(l[0])
	}
	return strings.Join(ls, ",")
}

func (c cjsChart) marshalSimpleData(col int) []string {
	ds := make([]string, len(c.inData.FSS))
	for i, f := range c.inData.FSS {
		ds[i] = fmt.Sprintf("%g", f[col])
	}
	return ds
}

func (c cjsChart) tooltipCallback() string {
	switch c.inData.ChartType {
	case "pie":
		return `
                    var value = data.datasets[0].data[tti.index];
                    var total = data.datasets[0].data.reduce((a, b) => a + b, 0)
                    var label = data.labels[tti.index];
                    var percentage = Math.round(value / total * 100);
                    return label + ': ' + percentage + '%';
    `
	case "line", "scatterline":
		return `
                    var value = data.datasets[tti.datasetIndex].data[tti.index];
                    if (value.y) {
                        value = value.y
                    }
                    return value;
    `
	case "scatter":
		return `
                    var value = data.datasets[tti.datasetIndex].data[tti.index];
                    var label = data.datasets[tti.datasetIndex].label;
                    return (label ? label + ': ' : '') + '(' + value.x + ', ' + value.y + ')';
    `
	case "bar":
		return `
                    var value = data.datasets[0].data[tti.index];
                    var label = data.labels[tti.index];
                    return value;
    `
	default:
		return ``
	}
}

func scatterRadius(x, min, max float64) float64 {
	if max-min < 50 {
		return x - min + 4
	}
	return float64(4) + (x-min)/(max-min)*50
}
