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
			name: "one column, all diferent",
			isss: [][]string{
				[]string{"a"},
				[]string{"b"},
				[]string{"c"},
			},
			expectedFSS: [][]float64{
				[]float64{1},
				[]float64{1},
				[]float64{1},
			},
			expectedSSS: [][]string{
				[]string{"a"},
				[]string{"b"},
				[]string{"c"},
			},
		},
		{
			name: "one column, one duplicate",
			isss: [][]string{
				[]string{"a"},
				[]string{"b"},
				[]string{"a"},
			},
			expectedFSS: [][]float64{
				[]float64{2},
				[]float64{1},
			},
			expectedSSS: [][]string{
				[]string{"a"},
				[]string{"b"},
			},
		},
		{
			name: "two columns, ignores second one",
			isss: [][]string{
				[]string{"a", "ignore me"},
				[]string{"b", "ignore me"},
				[]string{"a", "ignore me"},
			},
			expectedFSS: [][]float64{
				[]float64{2},
				[]float64{1},
			},
			expectedSSS: [][]string{
				[]string{"a"},
				[]string{"b"},
			},
		},
		{
			name: "sums less frequent labels into an Other category when over 10 labels",
			isss: [][]string{
				[]string{"a"},
				[]string{"a"},
				[]string{"b"},
				[]string{"b"},
				[]string{"c"},
				[]string{"c"},
				[]string{"d"},
				[]string{"d"},
				[]string{"e"},
				[]string{"e"},
				[]string{"f"},
				[]string{"f"},
				[]string{"g"},
				[]string{"g"},
				[]string{"h"},
				[]string{"h"},
				[]string{"i"},
				[]string{"i"},
				[]string{"j"},
				[]string{"k"},
				[]string{"l"},
				[]string{"m"},
			},
			expectedFSS: [][]float64{
				[]float64{2},
				[]float64{2},
				[]float64{2},
				[]float64{2},
				[]float64{2},
				[]float64{2},
				[]float64{2},
				[]float64{2},
				[]float64{2},
				[]float64{4},
			},
			expectedSSS: [][]string{
				[]string{"a"},
				[]string{"b"},
				[]string{"c"},
				[]string{"d"},
				[]string{"e"},
				[]string{"f"},
				[]string{"g"},
				[]string{"h"},
				[]string{"i"},
				[]string{"Other"},
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
