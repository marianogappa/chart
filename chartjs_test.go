package main

import "testing"

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
