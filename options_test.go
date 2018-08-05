package main

import (
	"reflect"
	"testing"

	"github.com/marianogappa/chart/format"
)

func TestResolveOptions(t *testing.T) {
	tests := []struct {
		args     []string
		expected options
		fails    bool
	}{
		{
			args: []string{},
			expected: options{
				title:      "",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t'},
			},
		},
		{
			args: []string{"-t", "title"},
			expected: options{
				title:      "title",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t'},
			},
		},
		{
			args: []string{"-title", "title"},
			expected: options{
				title:      "title",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t'},
			},
		},
		{
			args:  []string{"-t"},
			fails: true,
		},
		{
			args:  []string{"-title"},
			fails: true,
		},
		{
			args: []string{"bar", "log", "invert", ","},
			expected: options{
				title:      "",
				separator:  ',',
				scaleType:  logarithmic,
				chartType:  bar,
				lineFormat: format.LineFormat{Separator: ','},
			},
		},
		{
			args: []string{"bar", ";"},
			expected: options{
				title:      "",
				separator:  ';',
				chartType:  bar,
				lineFormat: format.LineFormat{Separator: ';'},
			},
		},
		{
			args: []string{" "},
			expected: options{
				title:      "",
				separator:  ' ',
				lineFormat: format.LineFormat{Separator: ' '},
			},
		},
		{
			args: []string{" ", "pie"},
			expected: options{
				title:      "",
				separator:  ' ',
				chartType:  pie,
				lineFormat: format.LineFormat{Separator: ' '},
			},
		},
		{
			args: []string{"line"},
			expected: options{
				title:      "",
				separator:  '\t',
				chartType:  line,
				lineFormat: format.LineFormat{Separator: '\t'},
			},
		},
		{
			args: []string{"scatter"},
			expected: options{
				title:      "",
				separator:  '\t',
				chartType:  scatter,
				lineFormat: format.LineFormat{Separator: '\t'},
			},
		},
		{
			args: []string{"-title", "title", "legacy-color", "1"},
			expected: options{
				title:      "title",
				separator:  '\t',
				colorType:  legacyColor,
				lineFormat: format.LineFormat{Separator: '\t'},
			},
		},
		{
			args: []string{"-title", "title", "gradient", "1"},
			expected: options{
				title:      "title",
				separator:  '\t',
				colorType:  gradient,
				lineFormat: format.LineFormat{Separator: '\t'},
			},
		},
		{
			args: []string{"ANSIC"},
			expected: options{
				dateFormat: "Mon Jan _2 15:04:05 2006",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t', DateFormat: "Mon Jan _2 15:04:05 2006"},
			},
		},
		{
			args: []string{"UnixDate"},
			expected: options{
				dateFormat: "Mon Jan _2 15:04:05 MST 2006",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t', DateFormat: "Mon Jan _2 15:04:05 MST 2006"},
			},
		},
		{
			args: []string{"RubyDate"},
			expected: options{
				dateFormat: "Mon Jan 02 15:04:05 -0700 2006",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t', DateFormat: "Mon Jan 02 15:04:05 -0700 2006"},
			},
		},
		{
			args: []string{"RFC822"},
			expected: options{
				dateFormat: "02 Jan 06 15:04 MST",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t', DateFormat: "02 Jan 06 15:04 MST"},
			},
		},
		{
			args: []string{"RFC822Z"},
			expected: options{
				dateFormat: "02 Jan 06 15:04 -0700",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t', DateFormat: "02 Jan 06 15:04 -0700"},
			},
		},
		{
			args: []string{"RFC850"},
			expected: options{
				dateFormat: "Monday, 02-Jan-06 15:04:05 MST",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t', DateFormat: "Monday, 02-Jan-06 15:04:05 MST"},
			},
		},
		{
			args: []string{"RFC1123"},
			expected: options{
				dateFormat: "Mon, 02 Jan 2006 15:04:05 MST",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t', DateFormat: "Mon, 02 Jan 2006 15:04:05 MST"},
			},
		},
		{
			args: []string{"RFC1123Z"},
			expected: options{
				dateFormat: "Mon, 02 Jan 2006 15:04:05 -0700",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t', DateFormat: "Mon, 02 Jan 2006 15:04:05 -0700"},
			},
		},
		{
			args: []string{"RFC3339"},
			expected: options{
				dateFormat: "2006-01-02T15:04:05Z07:00",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t', DateFormat: "2006-01-02T15:04:05Z07:00"},
			},
		},
		{
			args: []string{"RFC3339Nano"},
			expected: options{
				dateFormat: "2006-01-02T15:04:05.999999999Z07:00",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t', DateFormat: "2006-01-02T15:04:05.999999999Z07:00"},
			},
		},
		{
			args: []string{"Kitchen"},
			expected: options{
				dateFormat: "3:04PM",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t', DateFormat: "3:04PM"},
			},
		},
		{
			args: []string{"Stamp"},
			expected: options{
				dateFormat: "Jan _2 15:04:05",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t', DateFormat: "Jan _2 15:04:05"},
			},
		},
		{
			args: []string{"StampMilli"},
			expected: options{
				dateFormat: "Jan _2 15:04:05.000",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t', DateFormat: "Jan _2 15:04:05.000"},
			},
		},
		{
			args: []string{"StampMicro"},
			expected: options{
				dateFormat: "Jan _2 15:04:05.000000",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t', DateFormat: "Jan _2 15:04:05.000000"},
			},
		},
		{
			args: []string{"mysql"},
			expected: options{
				dateFormat: "2006-01-02 15:04:05",
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t', DateFormat: "2006-01-02 15:04:05"},
			},
		},
		{
			args: []string{"debug"},
			expected: options{
				separator:  '\t',
				lineFormat: format.LineFormat{Separator: '\t'},
				debug:      true,
			},
		},
		{
			args:  []string{"help"},
			fails: true,
		},
	}

	for _, ts := range tests {
		result, err := resolveOptions(ts.args)

		if ts.fails && err == nil {
			t.Errorf("should have failed with %v", result)
		}

		if !ts.fails && err != nil {
			t.Errorf("should not have failed resolving options with %v", result)
		}

		if !ts.fails && !reflect.DeepEqual(result, ts.expected) {
			t.Errorf("options are incorrect: %v was not equal to %v", result, ts.expected)
		}
	}
}
