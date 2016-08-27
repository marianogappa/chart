package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseLineFormat(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		sep      rune
		expected string
	}{
		{
			name:     "empty case",
			s:        "",
			sep:      '\t',
			expected: " ",
		},
		{
			name:     "basic tab",
			s:        "1\t2\t3",
			sep:      '\t',
			expected: "f,f,f",
		},
		{
			name:     "basic space",
			s:        "1 2 3",
			sep:      ' ',
			expected: "f,f,f",
		},
		{
			name:     "basic comma",
			s:        "1,2,3",
			sep:      ',',
			expected: "f,f,f",
		},
		{
			name:     "basic semicolon",
			s:        "1;2;3",
			sep:      ';',
			expected: "f,f,f",
		},
		{
			name:     "space with extras before, in between and after",
			s:        "  1   2  3  ",
			sep:      ' ',
			expected: "f,f,f",
		},
		{
			name:     "commas with complete floating numbers",
			s:        "-1,2.0e3,-3.239847E-1",
			sep:      ',',
			expected: "f,f,f",
		},
		{
			name:     "subsequent commas",
			s:        ",,",
			sep:      ',',
			expected: "f,f,f",
		},
		{
			name:     "string and float",
			s:        "a,1",
			sep:      ',',
			expected: "s,f",
		},
		{
			name:     "float and string",
			s:        "1,a",
			sep:      ',',
			expected: "f,s",
		},
		{
			name:     "float and string; ignore other separators",
			s:        "1,a;b c\td",
			sep:      ',',
			expected: "f,s",
		},
	}

	for _, ts := range tests {
		result := parseLineFormat(ts.s, ts.sep)
		if !reflect.DeepEqual(result, ts.expected) {
			t.Errorf("'%v' failed: %v was not equal to %v", ts.name, result, ts.expected)
		}
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		name     string
		i        []string
		sep      rune
		expected string
	}{
		{
			name:     "empty case",
			i:        []string{},
			sep:      '\t',
			expected: "",
		},
		{
			name:     "string, float",
			i:        []string{"a\t1", "b\t2", "c\t3"},
			sep:      '\t',
			expected: "s,f",
		},
		{
			name:     "string, float with one outlier (minority)",
			i:        []string{"a\t1", "onlystring", "c\t3"},
			sep:      '\t',
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

func TestParseLine(t *testing.T) {
	tests := []struct {
		name            string
		i               string
		sep             rune
		format          string
		expectedFloats  []float64
		expectedStrings []string
		fails           bool
	}{
		{
			name:            "base case",
			i:               "1,2,3",
			sep:             ',',
			format:          "f,f,f",
			expectedFloats:  []float64{1, 2, 3},
			expectedStrings: []string{},
			fails:           false,
		},
		{
			name:            "base failing case",
			i:               "1,a,3",
			sep:             ',',
			format:          "f,f,f",
			expectedFloats:  []float64{},
			expectedStrings: []string{},
			fails:           true,
		},
		{
			name:            "with strings",
			i:               "a,1",
			sep:             ',',
			format:          "s,f",
			expectedFloats:  []float64{1},
			expectedStrings: []string{"a"},
			fails:           false,
		},
		{
			name:            "strings and extra whitespace",
			i:               "    a   ,   1   ",
			sep:             ',',
			format:          "s,f",
			expectedFloats:  []float64{1},
			expectedStrings: []string{"a"},
			fails:           false,
		},
	}

	for _, ts := range tests {
		fs, ss, err := parseLine(ts.i, ts.format, ts.sep)
		if err != nil && !ts.fails {
			t.Errorf("'%v' failed: should not have failed but did! With [%v]", ts.name, err)
		}
		if err == nil && ts.fails {
			t.Errorf("'%v' failed: should have failed but didn't!", ts.name)
		}
		if !ts.fails && !reflect.DeepEqual(fs, ts.expectedFloats) {
			t.Errorf("'%v' failed: %v != %v", ts.name, fs, ts.expectedFloats)
		}
		if !ts.fails && !reflect.DeepEqual(ss, ts.expectedStrings) {
			fmt.Println(len(ss), len(ts.expectedStrings), ss == nil, ts.expectedStrings == nil)
			t.Errorf("'%v' failed: %v != %v", ts.name, ss, ts.expectedStrings)
		}
	}
}
