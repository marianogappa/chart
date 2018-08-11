package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"text/template"

	"github.com/marianogappa/chart/chartjs"
	"github.com/marianogappa/chart/format"
)

type cheatsheetExample struct {
	ID          int
	Title       string
	OptionsLine string
	Lines       []string
	HTML        string
}

func TestCheatsheet(t *testing.T) {
	fs, err := filepath.Glob("testdata/*.csv")
	if err != nil {
		t.Errorf("Could not get testdata [%v]", err)
		t.FailNow()
	}

	cheetsheetExamples := make([]cheatsheetExample, len(fs))
	for j, f := range fs {
		fh, err := os.Open(f)
		if err != nil {
			t.Errorf("[%v] Could not open file", f)
			t.FailNow()
		}

		rd := bufio.NewReader(fh)
		optionsLine, err := rd.ReadString('\n')
		if err == io.EOF {
			t.Errorf("[%v] doesn't have options line", f)
			t.FailNow()
		}

		r := regexp.MustCompile("'.+'|\".+\"|\\S+")
		m := r.FindAllString(optionsLine, -1)

		o, err := resolveOptions(m)
		log.Println(o.title)
		if err != nil {
			t.Errorf("[%v] can't resolve options line [%v]", f, optionsLine)
			t.FailNow()
		}
		if _, err = rd.ReadString('\n'); err != nil {
			t.Errorf("[%v] doesn't have empty line after options line", f)
			t.FailNow()
		}

		var rdr io.Reader
		rdr, o.lineFormat = format.Parse(rd, o.separator, o.dateFormat)
		d := mustNewDataset(rdr, o.lineFormat)
		if !o.lineFormat.HasFloats && !o.lineFormat.HasDateTimes && o.lineFormat.HasStrings {
			d.fss, d.sss, o.lineFormat = preprocessFreq(d.sss, o.lineFormat)
		}
		o.chartType, err = resolveChartType(o.chartType, o.lineFormat, d.Len())
		if err != nil {
			t.Errorf("[%v] error resolving chart type when o.chartType=%v and d.lineFormat=%v: [%v]", f, o.chartType, o.lineFormat, err)
			t.FailNow()
		}
		b, err := chartjs.New(
			o.chartType.String(),
			d.fss,
			d.sss,
			d.tss,
			d.minFSS,
			d.maxFSS,
			o.title,
			o.scaleType.String(),
			o.xLabel,
			o.yLabel,
			o.zeroBased,
			int(o.colorType),
		).Build()
		if err != nil {
			t.Errorf("[%v] breaks building chart with: [%v]", f, err)
			t.FailNow()
		}

		cheetsheetExamples[j] = cheatsheetExample{
			ID:          j,
			Title:       f,
			OptionsLine: optionsLine,
			Lines:       mustReadLines(f, t),
			HTML:        b.String(),
		}
	}

	cheetsheetExamplesTemplate, err := template.New("").Parse(cheetsheetExamplesTemplateString)
	if err != nil {
		t.Errorf("Could not parse cheatsheet examples page: [%v]", err)
		t.FailNow()
	}

	var b bytes.Buffer
	if err := cheetsheetExamplesTemplate.Execute(&b, cheetsheetExamples); err != nil {
		t.Errorf("[Could not execute HTML template for cheatsheet [%v]", err)
		t.FailNow()
	}

	err = ioutil.WriteFile("index.html", b.Bytes(), os.ModePerm)
	if err != nil {
		t.Errorf("Could not write file [%v] [%v]", "index.html", err)
		t.FailNow()
	}
}

func mustReadLines(path string, t *testing.T) []string {
	file, err := os.Open(path)
	if err != nil {
		t.Errorf("[%v] Could not open file", path)
		t.FailNow()
	}
	defer file.Close()

	var (
		lines []string
		i     int
	)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i++
		if i < 3 {
			continue
		}
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		t.Errorf("[%v] Error scanning file", scanner.Err())
		t.FailNow()
	}
	return lines
}

var cheetsheetExamplesTemplateString = baseTemplateHeaderString + cheetsheetString + baseTemplateFooterString
var cheetsheetString = `
<style>
body {
  background-color: #EEE;
  padding: 15px;
}
h1 {
  font-size: 30px;
  padding-top: 15px;
  padding-bottom: 25px;
  padding-left: 15px;
}
h2 {
  font-family: Lucida Console, Monaco, monospace;
  font-size: 13px;
  padding-bottom: 15px;
  color: #00FF00;
}
.chart {
  height: 300px;
  width: 650px;
  padding-bottom: 15px;
  background-color: #FFF;
  padding: 10px;
  padding-top: 30px;
  padding-bottom: 35px;
  border-radius: 10px;
}
.pre {
  font-family: Lucida Console, Monaco, monospace;
  font-size: 13px;
  display: block;
  padding-bottom: 15px;
  color: #00FF00;
}
.chart-wrapper {
  background-color: #000;
  border-radius: 5px;
  padding: 15px;
  margin: 10px;
  width: 670px;
  display: inline-block;
  box-shadow: #333 1px 1px 8px;
  vertical-align: top;
}

</style>
<h1>Chart Cheatsheet</h1>
{{range .}}
<div class="chart-wrapper">
  <h2><pre>$ cat {{.Title}}</pre></h2>
  <div class="pre"><pre>{{range .Lines}}{{.}}<br>{{end}}</pre></div>
  <h2><pre>$ cat {{.Title}} | chart {{.OptionsLine}}</pre></h2>
  <div class="chart">
    <canvas id="chart{{.ID}}" height="300px" width="650px"></canvas>
  </div>
  <script>
    var ctx{{.ID}} = document.getElementById("chart{{.ID}}");
    var chart{{.ID}} = new Chart(ctx{{.ID}}, {{.HTML}});
  </script>
</div>
{{end}}
`
