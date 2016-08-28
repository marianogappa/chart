package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestInput(t *testing.T) {
	tests := []struct {
		file          string
		expectedLines []string
	}{
		{
			file:          ``,
			expectedLines: []string{},
		},
		{
			file:          `a`,
			expectedLines: []string{},
		},
		{
			file: `a
`,
			expectedLines: []string{"a"},
		},
		{
			file: `   a
			`,
			expectedLines: []string{"a"},
		},
		{
			file: `   a
			b
			           c
			           `,
			expectedLines: []string{"a", "b", "c"},
		},
	}

	for _, ts := range tests {
		ls := readInput(strings.NewReader(ts.file))

		if !reflect.DeepEqual(ls, ts.expectedLines) {
			t.Errorf("options are incorrect: %v was not equal to %v", ls, ts.expectedLines)
		}
	}
}
