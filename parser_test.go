package curl2go

import (
	"reflect"
	"testing"
)

func TestParser_FlagParse(t *testing.T) {
	testCases := []struct {
		input    string
		expected ParsedFlags
	}{
		{
			input: "curl https://example.com",
			expected: ParsedFlags{
				UnFlags:      []string{"curl", "https://example.com"},
				BoolFlags:    map[string]bool{},
				StringsFlags: make(map[string][]string),
			},
		},
		{
			input: "curl -X POST https://example.com",
			expected: ParsedFlags{
				UnFlags: []string{"curl", "https://example.com"},
				StringsFlags: map[string][]string{
					"request": {"POST"},
				},
				BoolFlags: map[string]bool{},
			},
		},
		{
			input: "curl --request POST https://example.com",
			expected: ParsedFlags{
				UnFlags: []string{"curl", "https://example.com"},
				StringsFlags: map[string][]string{
					"request": {"POST"},
				},
				BoolFlags: map[string]bool{},
			},
		},
		{
			input: "curl -H 'Content-Type: application/json' https://example.com",
			expected: ParsedFlags{
				UnFlags: []string{"curl", "https://example.com"},
				StringsFlags: map[string][]string{
					"header": {"Content-Type: application/json"},
				},
				BoolFlags: map[string]bool{},
			},
		},
	}

	parser := NewParser()
	for _, tc := range testCases {
		actual := parser.FlagParse(tc.input)

		if !reflect.DeepEqual(actual, tc.expected) {
			t.Fatalf("FlagParse(%s) = %+v, expected %+v", tc.input, actual, tc.expected)
		}
	}
}
