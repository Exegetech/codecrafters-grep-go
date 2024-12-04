package main

import (
	"testing"
)

var tests = []struct {
	regex  string
	text   string
	result bool
}{
	{`a`, "apple", true},
	{`a`, "dog", false},

	{`\d`, "apple123", true},
	{`\d`, "123", true},
	{`\d`, "apple", false},

	{`\w`, "alpha-num3ric", true},

	{`[abc]`, "apple", true},
	{`[abc]`, "dog", false},

	{`[^xyz]`, "apple", true},
	{`[^anb]`, "banana", false},

	{`\d apple`, "1 apple", true},
	{`\d apple`, "3 oranges", false},
	{`\d\d\d apples`, "sally has 124 apples", true},
	{`\d\\d\\d apples`, "sally has 12 apples", false},
	{`\d \w\w\ws`, "sally has 3 dogs", true},
	{`\d \w\w\ws`, "sally has 4 dogs", true},
	{`\d \w\w\ws`, "sally has 1 dog", false},

	{`^log`, "log", true},
	{`^log`, "does not log", false},

	{`dog$`, "dog", true},
	{`dog$`, "dog is cute", false},

	{`ca+ts`, "cats", true},
	{`ca+ts`, "caats", true},
	{`ca+ts`, "cts", false},

	{`ca?t`, "cat", true},
	{`ca?t`, "act", true},
	{`ca?t`, "dog", false},
	{`ca?t`, "cag", false},

	{`d.g`, "dog", true},
	{`d.g`, "dg", false},
	{`g.+gol`, "goaoaoaoagol", true},
	{`g.+gol`, "gol", false},

	{`(cat|dog)`, "cat", true},
	{`(cat|dog)`, "dog", true},
	{`(cat|dog|man)`, "man", true},
	{`(cat|dog|man)`, "fish", false},
	{`a (cat|dog) and (cat|dog)s`, "a dog and cats", true},
	{`a (cat|dog)`, "a cow", false},
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
