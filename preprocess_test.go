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
		expected lineFormat
	}{
		{
			name: "empty case",
			s: "",
			sep: '\t',
			expected: lineFormat{emptyT},
		},
		{
			name: "basic tab",
			s: "1\t2\t3",
			sep: '\t',
			expected: lineFormat{floatOrStringT, separatorT, floatOrStringT, separatorT, floatOrStringT},
		},
		{
			name: "basic space",
			s: "1 2 3",
			sep: ' ',
			expected: lineFormat{floatOrStringT, separatorT, floatOrStringT, separatorT, floatOrStringT},
		},
		{
			name: "basic comma",
			s: "1,2,3",
			sep: ',',
			expected: lineFormat{floatOrStringT, separatorT, floatOrStringT, separatorT, floatOrStringT},
		},
		{
			name: "basic semicolon",
			s: "1;2;3",
			sep: ';',
			expected: lineFormat{floatOrStringT, separatorT, floatOrStringT, separatorT, floatOrStringT},
		},
		{
			name: "space with extras before, in between and after",
			s: "  1   2  3  ",
			sep: ' ',
			expected: lineFormat{floatOrStringT, separatorT, floatOrStringT, separatorT, floatOrStringT},
		},
		{
			name: "commas with complete floating numbers",
			s: "-1,2.0e3,-3.239847E-1",
			sep: ',',
			expected: lineFormat{floatOrStringT, separatorT, floatOrStringT, separatorT, floatOrStringT},
		},
		{
			name: "subsequent commas",
			s: ",,",
			sep: ',',
			expected: lineFormat{floatOrStringT, separatorT, floatOrStringT, separatorT, floatOrStringT},
		},
		{
			name: "string and float",
			s: "a,1",
			sep: ',',
			expected: lineFormat{stringT, separatorT, floatOrStringT},
		},
		{
			name: "float and string",
			s: "1,a",
			sep: ',',
			expected: lineFormat{floatOrStringT, separatorT, stringT},
		},
		{
			name: "float and string; ignore other separators",
			s: "1,a;b c\td",
			sep: ',',
			expected: lineFormat{floatOrStringT, separatorT, stringT},
		},
	}

	for _, ts := range tests {
		result := parseLine(ts.s, ts.sep)
		if !reflect.DeepEqual(result, ts.expected) {
			t.Errorf("'%v' failed: %v was not equal to %v", ts.name, result, ts.expected)
		}
	}
}
