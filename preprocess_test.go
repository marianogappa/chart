package main

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name string
		s     string
		sep rune
		expected string
	}{
		{
			name: "empty case",
			s: "",
			sep: '\t',
			expected: " ",
		},
		{
			name: "basic tab",
			s: "1\t2\t3",
			sep: '\t',
			expected: "f,f,f",
		},
		{
			name: "basic space",
			s: "1 2 3",
			sep: ' ',
			expected: "f,f,f",
		},
		{
			name: "basic comma",
			s: "1,2,3",
			sep: ',',
			expected: "f,f,f",
		},
		{
			name: "basic semicolon",
			s: "1;2;3",
			sep: ';',
			expected: "f,f,f",
		},
		{
			name: "space with extras before, in between and after",
			s: "  1   2  3  ",
			sep: ' ',
			expected: "f,f,f",
		},
		{
			name: "commas with complete floating numbers",
			s: "-1,2.0e3,-3.239847E-1",
			sep: ',',
			expected: "f,f,f",
		},
		{
			name: "subsequent commas",
			s: ",,",
			sep: ',',
			expected: "f,f,f",
		},
		{
			name: "string and float",
			s: "a,1",
			sep: ',',
			expected: "s,f",
		},
		{
			name: "float and string",
			s: "1,a",
			sep: ',',
			expected: "f,s",
		},
		{
			name: "float and string; ignore other separators",
			s: "1,a;b c\td",
			sep: ',',
			expected: "f,s",
		},
	}

	for _, ts := range tests {
		result := parseLine(ts.s, ts.sep)
		if !reflect.DeepEqual(result, ts.expected) {
			t.Errorf("'%v' failed: %v was not equal to %v", ts.name, result, ts.expected)
		}
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		name string
		i    []string
		sep rune
		expected string
	}{
		{
			name: "empty case",
			i: []string{},
			sep: '\t',
			expected: "",
		},
		{
			name: "string, float",
			i: []string{"a\t1", "b\t2", "c\t3"},
			sep: '\t',
			expected: "s,f",
		},
		{
			name: "string, float with one outlier (minority)",
			i: []string{"a\t1", "onlystring", "c\t3"},
			sep: '\t',
			expected: "s,f",
		},
	}

	for _, ts := range tests {
		result := parseFormat(ts.i, ts.sep)
		if !reflect.DeepEqual(result, ts.expected) {
			t.Errorf("'%v' failed: %v was not equal to %v", ts.name, result, ts.expected)
		}
	}
}
