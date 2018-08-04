package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseLineFormat(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		sep      rune
		df       string
		expected string
	}{
		{
			name:     "empty case",
			s:        "",
			sep:      '\t',
			expected: "s",
		},
		{
			name:     "basic tab",
			s:        "1\t2\t3",
			sep:      '\t',
			expected: "fff",
		},
		{
			name:     "basic space",
			s:        "1 2 3",
			sep:      ' ',
			expected: "fff",
		},
		{
			name:     "basic comma",
			s:        "1,2,3",
			sep:      ',',
			expected: "fff",
		},
		{
			name:     "basic semicolon",
			s:        "1;2;3",
			sep:      ';',
			expected: "fff",
		},
		{
			name:     "space with extras before, in between and after",
			s:        "  1   2  3  ",
			sep:      ' ',
			expected: "fff",
		},
		{
			name:     "commas with complete floating numbers",
			s:        "-1,2.0e3,-3.239847E-1",
			sep:      ',',
			expected: "fff",
		},
		{
			name:     "subsequent commas",
			s:        ",,",
			sep:      ',',
			expected: "ss",
		},
		{
			name:     "string and float",
			s:        "a,1",
			sep:      ',',
			expected: "sf",
		},
		{
			name:     "float and string",
			s:        "1,a",
			sep:      ',',
			expected: "fs",
		},
		{
			name:     "float and string; ignore other separators",
			s:        "1,a;b c\td",
			sep:      ',',
			expected: "fs",
		},
		{
			name: "date on the left",
			s: "2016-08-29	0.0125", // N.B. hello; I've broken gofmt :)
			sep:      '\t',
			df:       "2006-01-02",
			expected: "df",
		},
		{
			name: "date on the right",
			s: "0.0125	2016-08-29",
			sep:      '\t',
			df:       "2006-01-02",
			expected: "fd",
		},
	}

	for _, ts := range tests {
		result := parseLineFormat(ts.s, ts.sep, ts.df)
		if !reflect.DeepEqual(result, ts.expected) {
			t.Errorf("'%v' failed: %v was not equal to %v", ts.name, result, ts.expected)
		}
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		name     string
		i        string
		sep      rune
		df       string
		expected string
	}{
		{
			name:     "empty case",
			i:        ``,
			sep:      '\t',
			expected: "",
		},
		{
			name: "string, float",
			i: `a	1
b	2
c	3
`,
			sep:      '\t',
			expected: "sf",
		},
		{
			name: "string, float with one outlier (minority)",
			i: `a	1
onlystring
c	3
`,
			sep:      '\t',
			expected: "sf",
		},
	}

	for _, ts := range tests {
		_, result := readAndParseFormat(strings.NewReader(ts.i), ts.sep, ts.df)
		if !reflect.DeepEqual(result, ts.expected) {
			t.Errorf("'%v' failed: %v was not equal to %v", ts.name, result, ts.expected)
		}
	}
}
