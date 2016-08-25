package main

import (
	"reflect"
	"testing"
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
				title:     "",
				separator: "\t",
				scaleType: linear,
				invert:    false,
				chartType: pie,
			},
		},
		{
			args: []string{"-t", "title"},
			expected: options{
				title:     "title",
				separator: "\t",
				scaleType: linear,
				invert:    false,
				chartType: pie,
			},
		},
		{
			args: []string{"-title", "title"},
			expected: options{
				title:     "title",
				separator: "\t",
				scaleType: linear,
				invert:    false,
				chartType: pie,
			},
		},
		{
			args: []string{"-separator", "\t"},
			expected: options{
				title:     "",
				separator: "\t",
				scaleType: linear,
				invert:    false,
				chartType: pie,
			},
		},
		{
			args: []string{"-separator", " "},
			expected: options{
				title:     "",
				separator: " ",
				scaleType: linear,
				invert:    false,
				chartType: pie,
			},
		},
		{
			args: []string{"-s", " "},
			expected: options{
				title:     "",
				separator: " ",
				scaleType: linear,
				invert:    false,
				chartType: pie,
			},
		},
		{
			args: []string{"-separator", ";"},
			expected: options{
				title:     "",
				separator: ";",
				scaleType: linear,
				invert:    false,
				chartType: pie,
			},
		},
		{
			args: []string{"-separator", ","},
			expected: options{
				title:     "",
				separator: ",",
				scaleType: linear,
				invert:    false,
				chartType: pie,
			},
		},
		{
			args: []string{"-separator", "invalid"},
			expected: options{
				title:     "",
				separator: "\t",
				scaleType: linear,
				invert:    false,
				chartType: pie,
			},
		},
		{
			args: []string{"-l"},
			expected: options{
				title:     "",
				separator: "\t",
				scaleType: logarithmic,
				invert:    false,
				chartType: pie,
			},
		},
		{
			args: []string{"-log"},
			expected: options{
				title:     "",
				separator: "\t",
				scaleType: logarithmic,
				invert:    false,
				chartType: pie,
			},
		},
		{
			args: []string{"-i"},
			expected: options{
				title:     "",
				separator: "\t",
				scaleType: linear,
				invert:    true,
				chartType: pie,
			},
		},
		{
			args: []string{"-invert"},
			expected: options{
				title:     "",
				separator: "\t",
				scaleType: linear,
				invert:    true,
				chartType: pie,
			},
		},
		{
			args: []string{"-y", "bar"},
			expected: options{
				title:     "",
				separator: "\t",
				scaleType: linear,
				invert:    false,
				chartType: bar,
			},
		},
		{
			args: []string{"-type", "bar"},
			expected: options{
				title:     "",
				separator: "\t",
				scaleType: linear,
				invert:    false,
				chartType: bar,
			},
		},
		{
			args: []string{"-type"},
			fails: true,
		},
		{
			args: []string{"-y"},
			fails: true,
		},
		{
			args: []string{"-t"},
			fails: true,
		},
		{
			args: []string{"-title"},
			fails: true,
		},
		{
			args: []string{"-separator"},
			fails: true,
		},
		{
			args: []string{"-s"},
			fails: true,
		},
		{
			args: []string{"bar", "log", "invert", ","},
			expected: options{
				title:     "",
				separator: ",",
				scaleType: logarithmic,
				invert:    true,
				chartType: bar,
			},
		},
		{
			args: []string{"bar", ";"},
			expected: options{
				title:     "",
				separator: ";",
				scaleType: linear,
				invert:    false,
				chartType: bar,
			},
		},
		{
			args: []string{" "},
			expected: options{
				title:     "",
				separator: " ",
				scaleType: linear,
				invert:    false,
				chartType: pie,
			},
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
