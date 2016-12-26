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
			t.Errorf("[%v] doesn't have options line", fs)
			t.FailNow()
		}

		r := regexp.MustCompile("'.+'|\".+\"|\\S+")
		m := r.FindAllString(optionsLine, -1)

		o, err := resolveOptions(m)
		log.Println(o.title)
		if err != nil {
			t.Errorf("[%v] can't resolve options line [%v]", fs, optionsLine)
			t.FailNow()
		}
		if _, err = rd.ReadString('\n'); err != nil {
			t.Errorf("[%v] doesn't have empty line after options line", fs)
			t.FailNow()
		}

		i := readInput(rd)
		if len(i) == 0 {
			t.Errorf("[%v] doesn't have data", fs)
			t.FailNow()
		}
		defer fh.Close()

		b, err := buildChart(i, o)
		if err != nil {
			t.Errorf("[%v] breaks building chart with: [%v]", fs, err)
			t.FailNow()
		}

		cheetsheetExamples[j] = cheatsheetExample{ID: j, Title: f, OptionsLine: optionsLine, Lines: i, HTML: b.String()}
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

	err = ioutil.WriteFile("cheatsheet.html", b.Bytes(), os.ModePerm)
	if err != nil {
		t.Errorf("Could not write file [%v] [%v]", "cheatsheet.html", err)
		t.FailNow()
	}
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
    <canvas id="chart{{.ID}}"></canvas>
  </div>
  <script>
    var ctx{{.ID}} = document.getElementById("chart{{.ID}}");
    var chart{{.ID}} = new Chart(ctx{{.ID}}, {{.HTML}});
  </script>
</div>
{{end}}
`
