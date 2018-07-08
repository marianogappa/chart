package main

import "testing"

func TestFreq(t *testing.T) {
	tests := []struct {
		name        string
		isss        [][]string
		expectedFSS [][]float64
		expectedSSS [][]string
	}{
		{
			name:        "empty case",
			isss:        [][]string{},
			expectedFSS: [][]float64{},
			expectedSSS: [][]string{},
		},
		{
			name: "one column, all different",
			isss: [][]string{
				{"a"},
				{"b"},
				{"c"},
			},
			expectedFSS: [][]float64{
				{1},
				{1},
				{1},
			},
			expectedSSS: [][]string{
				{"a"},
				{"b"},
				{"c"},
			},
		},
		{
			name: "one column, one duplicate",
			isss: [][]string{
				{"a"},
				{"b"},
				{"a"},
			},
			expectedFSS: [][]float64{
				{2},
				{1},
			},
			expectedSSS: [][]string{
				{"a"},
				{"b"},
			},
		},
		{
			name: "two columns, ignores second one",
			isss: [][]string{
				{"a", "ignore me"},
				{"b", "ignore me"},
				{"a", "ignore me"},
			},
			expectedFSS: [][]float64{
				{2},
				{1},
			},
			expectedSSS: [][]string{
				{"a"},
				{"b"},
			},
		},
		{
			name: "sums less frequent labels into an Other category when over 10 labels",
			isss: [][]string{
				{"a"},
				{"a"},
				{"b"},
				{"b"},
				{"c"},
				{"c"},
				{"d"},
				{"d"},
				{"e"},
				{"e"},
				{"f"},
				{"f"},
				{"g"},
				{"g"},
				{"h"},
				{"h"},
				{"i"},
				{"i"},
				{"j"},
				{"k"},
				{"l"},
				{"m"},
			},
			expectedFSS: [][]float64{
				{2},
				{2},
				{2},
				{2},
				{2},
				{2},
				{2},
				{2},
				{2},
				{4},
			},
			expectedSSS: [][]string{
				{"a"},
				{"b"},
				{"c"},
				{"d"},
				{"e"},
				{"f"},
				{"g"},
				{"h"},
				{"i"},
				{"Other"},
			},
		},
	}

	for _, ts := range tests {
		fss, sss := preprocessFreq(ts.isss)

		if !equalMap(fss, ts.expectedFSS, sss, ts.expectedSSS) {
			t.Errorf("[%v] case failed: %v, %v were not equal to %v, %v", ts.name, fss, sss, ts.expectedFSS, ts.expectedSSS)
		}
	}
}

// equalMap is different to intuitive equality in that: supports mirror slices sss/fss and disregards order of map entries
func equalMap(fss [][]float64, expectedFSS [][]float64, sss [][]string, expectedSSS [][]string) bool {
	if len(fss) != len(expectedFSS) || len(sss) != len(expectedSSS) {
		return false
	}

	for i := range fss {
		f := fss[i][0]
		s := sss[i][0]

		k := -1
		for j := range expectedSSS {
			if expectedSSS[j][0] == s {
				k = j
				break
			}
		}

		if k == -1 {
			return false
		}

		if expectedFSS[k][0] != f {
			return false
		}
	}

	return true
}
