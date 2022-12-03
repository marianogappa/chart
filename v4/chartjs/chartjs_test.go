package chartjs

import (
	"reflect"
	"testing"
)

func TestCalculateMinMaxFSS(t *testing.T) {
	tests := []struct {
		fss    [][]float64
		minFSS []float64
		maxFSS []float64
	}{
		{
			fss:    [][]float64{},
			minFSS: nil,
			maxFSS: nil,
		},
		{
			fss: [][]float64{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			minFSS: []float64{1, 2, 3},
			maxFSS: []float64{7, 8, 9},
		},
		{
			fss: [][]float64{
				{1, 8, 3},
				{4, 2, 9},
				{7, 6, 5},
			},
			minFSS: []float64{1, 2, 3},
			maxFSS: []float64{7, 8, 9},
		},
		{
			fss: [][]float64{
				{1.2, 5.6},
				{7.8, 3.4},
			},
			minFSS: []float64{1.2, 3.4},
			maxFSS: []float64{7.8, 5.6},
		},
	}
	for _, tc := range tests {
		actualMinFSS, actualMaxFSS := calculateMinMaxFSS(tc.fss)
		if !reflect.DeepEqual(tc.minFSS, actualMinFSS) {
			t.Errorf("Expected %v but got %v", tc.minFSS, actualMinFSS)
		}
		if !reflect.DeepEqual(tc.maxFSS, actualMaxFSS) {
			t.Errorf("Expected %v but got %v", tc.maxFSS, actualMaxFSS)
		}
	}
}

func TestCropLongLabels(t *testing.T) {
	tests := []struct {
		s        string
		expected string
	}{
		{
			s:        "",
			expected: "",
		},
		{
			s:        "01234567890123456789012345678901234567890123456789",
			expected: "01234567890123456789012345678901234567890123456789",
		},
		{
			s:        "012345678901234567890123456789012345678901234567890",
			expected: "01234567890123456789012345678901234567890123456...",
		},
	}

	for _, ts := range tests {
		result := cropLongLabels(ts.s)

		if result != ts.expected {
			t.Errorf("cropping long labels: %v was not equal to %v", result, ts.expected)
		}
	}
}

func TestPreprocessLabel(t *testing.T) {
	tests := []struct {
		s        string
		expected string
	}{
		{
			s:        "",
			expected: "",
		},
		{
			s:        "012345678901234567890123456789012345678901234567890",
			expected: "`01234567890123456789012345678901234567890123456...`",
		},
		{
			s:        "012345678901234567890123456789012345`78901234567890",
			expected: "`012345678901234567890123456789012345\\`7890123456...`",
		},
		{
			s:        "hello\\",
			expected: "`hello\\\\`",
		},
		{
			s:        "he${llo",
			expected: "`he\\${llo`",
		},
	}

	for _, ts := range tests {
		result := preprocessLabel(ts.s)

		if result != ts.expected {
			t.Errorf("preprocessing labels: %v was not equal to %v", result, ts.expected)
		}
	}
}
