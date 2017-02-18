package main

import (
	"testing"
	"time"
)

func TestBar(t *testing.T) {
	tests := []struct {
		name      string
		fss       [][]float64
		sss       [][]string
		tss       [][]time.Time
		title     string
		scaleType scaleType
		fails     bool
	}{
		{
			name:  "empty case; should fail",
			fails: true,
		},
		{
			name:  "inconsistent number of values between data points and labels",
			fss:   [][]float64{[]float64{1}},
			fails: true,
		},
		{
			name:  "basic working example",
			fss:   [][]float64{[]float64{1}, []float64{2}, []float64{3}},
			sss:   [][]string{[]string{"a"}, []string{"b"}, []string{"c"}},
			title: "Basic working example",
		},
		{
			name:  "should use time column as discrete dimension",
			fss:   [][]float64{[]float64{1}, []float64{2}, []float64{3}},
			tss:   [][]time.Time{[]time.Time{time.Now()}, []time.Time{time.Now().Add(1 * time.Hour)}, []time.Time{time.Now().Add(2 * time.Hour)}},
			fails: false,
		},
		{
			name:  "should use first float column as discrete dimension", // https://github.com/marianogappa/chart/issues/7
			fss:   [][]float64{[]float64{0, 1}, []float64{2, 463008}},
			fails: false,
		},
	}

	for _, ts := range tests {
		templateData, resultBarTemplate, err := setupBar(ts.fss, ts.sss, ts.tss, ts.title, linear, "", "", false)
		if ts.fails && err == nil {
			t.Errorf("'%v' should have failed", ts.name)
		}

		if !ts.fails {
			if err != nil {
				t.Errorf("'%v' shouldn't have failed, but did with %v", ts.name, err)
			}
			if resultBarTemplate != barTemplate {
				t.Errorf("'%v' appears to not be using the hardcoded barTemplate", ts.name)
			}
			if templateData.(barTemplateData).ChartType != "bar" {
				t.Errorf("'%v' appears to not be returning a bar chart", ts.name)
			}
			if templateData.(barTemplateData).Title != ts.title {
				t.Errorf("'%v' did not use the specified title", ts.name)
			}
			ds := templateData.(barTemplateData).Data
			ss := templateData.(barTemplateData).Labels
			if len(ts.fss) != len(ds) {
				t.Errorf("'%v' is using a different number of data points (%v) than specified (%v)", ts.name, len(ds), len(ts.fss))
			}
			if len(ds) != len(ss) {
				t.Errorf("'%v' is returning %v data points, but %v labels", ts.name, len(ds), len(ss))
			}
		}
	}
}
