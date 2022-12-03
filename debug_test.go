package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/marianogappa/chart/dataset"
)

func TestDebug(t *testing.T) {
	var ts = []struct {
		name     string
		d        dataset.Dataset
		o        options
		err      error
		expected string
	}{
		{
			name: "base case",
			d:    dataset.Dataset{},
			o:    options{},
			err:  nil,
			expected: `Lines read	0
Line format inferred	
Lines used	0
Float column count	0
String column count	0
Date/Time column count	0
Chart type	???
Scale type	linear
`,
		},
		{
			name: "reports lines read",
			d:    dataset.Dataset{StdinLen: 123},
			o:    options{},
			err:  nil,
			expected: `Lines read	123
Line format inferred	
Lines used	0
Float column count	0
String column count	0
Date/Time column count	0
Chart type	???
Scale type	linear
`,
		},
		{
			name: "reports pie chart type",
			d:    dataset.Dataset{},
			o:    options{chartType: pie},
			err:  nil,
			expected: `Lines read	0
Line format inferred	
Lines used	0
Float column count	0
String column count	0
Date/Time column count	0
Chart type	pie
Scale type	linear
`,
		},
		{
			name: "reports line chart type",
			d:    dataset.Dataset{},
			o:    options{chartType: line},
			err:  nil,
			expected: `Lines read	0
Line format inferred	
Lines used	0
Float column count	0
String column count	0
Date/Time column count	0
Chart type	line
Scale type	linear
`,
		},
		{
			name: "reports scatter chart type",
			d:    dataset.Dataset{},
			o:    options{chartType: scatter},
			err:  nil,
			expected: `Lines read	0
Line format inferred	
Lines used	0
Float column count	0
String column count	0
Date/Time column count	0
Chart type	scatter
Scale type	linear
`,
		},
		{
			name: "reports bar chart type",
			d:    dataset.Dataset{},
			o:    options{chartType: bar},
			err:  nil,
			expected: `Lines read	0
Line format inferred	
Lines used	0
Float column count	0
String column count	0
Date/Time column count	0
Chart type	bar
Scale type	linear
`,
		},
		{
			name: "reports linear scaleType",
			d:    dataset.Dataset{},
			o:    options{scaleType: linear},
			err:  nil,
			expected: `Lines read	0
Line format inferred	
Lines used	0
Float column count	0
String column count	0
Date/Time column count	0
Chart type	???
Scale type	linear
`,
		},
		{
			name: "reports logarithmic scaleType",
			d:    dataset.Dataset{},
			o:    options{scaleType: logarithmic},
			err:  nil,
			expected: `Lines read	0
Line format inferred	
Lines used	0
Float column count	0
String column count	0
Date/Time column count	0
Chart type	???
Scale type	logarithmic
`,
		},
		{
			name: "reports tab separator",
			d:    dataset.Dataset{},
			o:    options{separator: '\t'},
			err:  nil,
			expected: `Lines read	0
Line format inferred	
Lines used	0
Float column count	0
String column count	0
Date/Time column count	0
Chart type	???
Scale type	linear
Separator	[tab]
`,
		},
		{
			name: "reports space separator",
			d:    dataset.Dataset{},
			o:    options{separator: ' '},
			err:  nil,
			expected: `Lines read	0
Line format inferred	
Lines used	0
Float column count	0
String column count	0
Date/Time column count	0
Chart type	???
Scale type	linear
Separator	[space]
`,
		},
		{
			name: "reports comma separator",
			d:    dataset.Dataset{},
			o:    options{separator: ','},
			err:  nil,
			expected: `Lines read	0
Line format inferred	
Lines used	0
Float column count	0
String column count	0
Date/Time column count	0
Chart type	???
Scale type	linear
Separator	[comma]
`,
		},
		{
			name: "reports semicolon separator",
			d:    dataset.Dataset{},
			o:    options{separator: ';'},
			err:  nil,
			expected: `Lines read	0
Line format inferred	
Lines used	0
Float column count	0
String column count	0
Date/Time column count	0
Chart type	???
Scale type	linear
Separator	[semicolon]
`,
		},
		{
			name: "reports title and axis labels",
			d:    dataset.Dataset{},
			o:    options{title: "sample title", xLabel: "sample x-label", yLabel: "sample y-label"},
			err:  nil,
			expected: `Lines read	0
Line format inferred	
Lines used	0
Float column count	0
String column count	0
Date/Time column count	0
Chart title	sample title
Chart horizontal axis label	sample x-label
Chart vertical axis label	sample y-label
Chart type	???
Scale type	linear
`,
		},
		{
			name: "reports dateFormat",
			d:    dataset.Dataset{},
			o:    options{dateFormat: "2006-01-02 15:04:05"},
			err:  nil,
			expected: `Lines read	0
Line format inferred	
Lines used	0
Float column count	0
String column count	0
Date/Time column count	0
Date format	2006-01-02 15:04:05
Chart type	???
Scale type	linear
`,
		},
		{
			name: "reports lines used",
			d:    dataset.Dataset{FSS: [][]float64{{}}, SSS: [][]string{{}}, TSS: [][]time.Time{{}}},
			o:    options{},
			err:  nil,
			expected: `Lines read	0
Line format inferred	
Lines used	1
Float column count	0
String column count	0
Date/Time column count	0
Chart type	???
Scale type	linear
`,
		},
		{
			name: "reports errors",
			d:    dataset.Dataset{},
			o:    options{scaleType: linear},
			err:  fmt.Errorf("sample error"),
			expected: `Error trying to chart: sample error
Lines read	0
Line format inferred	
Lines used	0
Float column count	0
String column count	0
Date/Time column count	0
Chart type	???
Scale type	linear
`,
		},
	}

	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			var a = renderDebug(tc.d, tc.o, tc.err)
			if tc.expected != a {
				t.Errorf("Expected [%v] but got [%v]", tc.expected, a)
			}
		})
	}
}
