package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/skratchdot/open-golang/open"
)

var baseTemplate, pieTemplate *template.Template

type pieTemplateData struct {
	Labels       string
	DisplayTitle bool
	Title        string
	ChartType    string
	Data         string
}

func main() {
	chartType := mustParseChartType()
	title, displayTitle := mustParseTitle()
	// defer func() { time.Sleep(5 * time.Second); os.Remove(tmpfile.Name()) }()

	var finalString string
	var err error
	switch chartType {
	case "pie":
		finalString, err = setupPie(title, displayTitle, '\t')
	}
	if err != nil {
		log.Fatal(err)
	}

	tmpfile, err := ioutil.TempFile("", "chartData")
	if err != nil {
		log.Fatal(err)
	}

	if err = baseTemplate.Execute(tmpfile, finalString); err != nil {
		log.Fatal(err)
	}
	if err = tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	open.Run("file://" + tmpfile.Name())
}

func setupPie(title string, displayTitle bool, separator rune) (string, error) {
	r := csv.NewReader(os.Stdin)
	r.Comma = separator
	r.Comment = '#'

	var labels []string
	var data []float64
	var d string

	for {
		var fS, s string
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		if len(record) == 0 {
			break
		}
		if len(record) == 1 {
			fS = record[0]
			s = ""
		}
		if len(record) >= 2 {
			s = record[0]
			fS = record[1]
		}
		f, err := strconv.ParseFloat(fS, 64)
		if err != nil {
			log.Printf("Ignoring this as it's not a number: [%v]", d)
			continue
		}
		data = append(data, f)
		labels = append(labels, `"`+s+`"`)
	}

	var ds []string
	for _, v := range data {
		ds = append(ds, strconv.FormatFloat(v, 'f', -1, 64))
	}
	stringData := strings.Join(ds, ",")
	stringLabels := strings.Join(labels, ",")

	templateData := pieTemplateData{
		ChartType:    "pie",
		Data:         stringData,
		Labels:       stringLabels,
		Title:        title,
		DisplayTitle: displayTitle,
	}

	var b bytes.Buffer
	if err := pieTemplate.Execute(&b, templateData); err != nil {
		log.Fatal(err)
	}

	return b.String(), nil
}

func mustParseChartType() string {
	if len(os.Args) < 2 {
		return "pie"
	}

	switch os.Args[1] {
	case "pie":
		return "pie"
	default:
		return "pie"
	}
}

func mustParseTitle() (string, bool) {
	if len(os.Args) < 3 {
		return "", false
	}
	return os.Args[2], true
}
