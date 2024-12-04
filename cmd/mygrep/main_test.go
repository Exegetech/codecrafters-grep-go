package main

import (
	"testing"
)

var tests = []struct {
	regex  string
	text   string
	result bool
}{
	{"a", "apple", true},
	{`\d`, "apple123", true},
}

func TestFlagParser(t *testing.T) {
	for _, tt := range tests {
		t.Run("Regex", func(t *testing.T) {
			result := match(tt.regex, tt.text)

			if result != tt.result {
				t.Errorf("regex %q with text %q, got %t, want %t", tt.regex, tt.text, result, tt.result)
			}
		})
	}
}
