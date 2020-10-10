package main

import (
	"strings"
	"testing"
)

func TestNumToHumanDecimal(t *testing.T) {
	var tests = []struct {
		input  string
		result string
	}{
		{"999", "999"},
		{"1000", "1K"},
		{"1500", "2K"},
		{"12345", "12K"},
		{"123456", "123K"},
		{"987654", "988K"},
		{"1234567", "1M"},
		{"12345678", "12M"},
		{"123456789", "123M"},
		{"123456789012", "123G"},
		{"123456789012345", "123T"},
		{"123456789012345678", "123P"},
		{"123456789012345678901", "123E"},
		{"123456789012345678901234", "123Z"},
		{"123456789012345678901234567", "123Y"},
		{"1234567890123456789012345678", "1235Y"}, // overflow our suffixes
		{"12345678901234567890123456789012345678901234567890", "12345678901234567890123456789012345678901234567890"}, // overflow float64
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			res := numToHuman(false, 1000)(test.input)
			if res != test.result {
				t.Errorf("got %s, want %s", res, test.result)
			}
		})
	}
}

func TestNumToHumanBinary(t *testing.T) {
	var tests = []struct {
		input  string
		result string
	}{
		{"999", "999"},
		{"1000", "1Ki"},
		{"1500", "1Ki"},
		{"1600", "2Ki"},
		{"12345", "12Ki"},
		{"123456", "121Ki"},
		{"987654", "965Ki"},
		{"1234567", "1Mi"},
		{"12345678", "12Mi"},
		{"123456789", "118Mi"},
		{"123456789012", "115Gi"},
		{"123456789012345", "112Ti"},
		{"123456789012345678", "110Pi"},
		{"123456789012345678901", "107Ei"},
		{"123456789012345678901234", "105Zi"},
		{"123456789012345678901234567", "102Yi"},
		{"1234567890123456789012345678", "1021Yi"}, // overflow our suffixes
		{"12345678901234567890123456789012345678901234567890", "12345678901234567890123456789012345678901234567890"}, // overflow float64
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			res := numToHuman(true, 1000)(test.input)
			if res != test.result {
				t.Errorf("got %s, want %s", res, test.result)
			}
		})
	}
}

func TestHumanize(t *testing.T) {
	var tests = []struct {
		input  string
		result string
	}{
		{
			"Here are some biggish numbers: 123456:7890123. Thanks!",
			"Here are some biggish numbers: 123K:8M. Thanks!",
		},
		{
			"This is too small: 6. And this one: 42.",
			"This is too small: 6. And this one: 42.",
		},
		{
			"And this doesn't even have any numbers.",
			"And this doesn't even have any numbers.",
		},
		{
			`This has
multiple lines, some with numbers like 987654321 and
some without.`,
			`This has
multiple lines, some with numbers like 988M and
some without.`,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			reader := strings.NewReader(test.input)
			var result strings.Builder
			humanize(reader, &result, false, 1000)
			if result.String() != test.result {
				t.Errorf("got %s, want %s", result.String(), test.result)
			}
		})
	}
}
