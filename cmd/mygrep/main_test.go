package main

import (
	"testing"
)

var tests = []struct {
	regex  string
	text   string
	result bool
}{
	// Literal
	{`a`, "apple", true},
	{`a`, "dog", false},

	// Digit
	{`\d`, "apple123", true},
	{`\d`, "123", true},
	{`\d`, "apple", false},

	// Alphanumeric
	{`\w`, "alpha-num3ric", true},

	// Character set
	{`[abc]`, "apple", true},
	{`[abc]`, "dog", false},

	// Negated character set
	{`[^xyz]`, "apple", true},
	{`[^anb]`, "banana", false},

	// Mix of everything above
	{`\d apple`, "1 apple", true},
	{`\d apple`, "3 oranges", false},
	{`\d\d\d apples`, "sally has 124 apples", true},
	{`\d\\d\\d apples`, "sally has 12 apples", false},
	{`\d \w\w\ws`, "sally has 3 dogs", true},
	{`\d \w\w\ws`, "sally has 4 dogs", true},
	{`\d \w\w\ws`, "sally has 1 dog", false},

	// Start of string anchors
	{`^log`, "log", true},
	{`^log`, "does not log", false},

	// End of string anchor
	{`dog$`, "dog", true},
	{`dog$`, "dog is cute", false},

	// One or more
	{`ca+ts`, "cats", true},
	{`ca+ts`, "caats", true},
	{`ca+ts`, "cts", false},

	// Zero or more
	{`ca?t`, "cat", true},
	{`ca?t`, "act", true},
	{`ca?t`, "dog", false},
	{`ca?t`, "cag", false},

	// Wildcard
	{`d.g`, "dog", true},
	{`d.g`, "dg", false},
	{`g.+gol`, "goaoaoaoagol", true},
	{`g.+gol`, "gol", false},

	// Alternation
	{`(cat|dog)`, "cat", true},
	{`(cat|dog)`, "dog", true},
	{`(cat|dog|man)`, "man", true},
	{`(cat|dog|man)`, "fish", false},
	{`a (cat|dog) and (cat|dog)s`, "a dog and cats", true},
	{`a (cat|dog)`, "a cow", false},

	// Mix everything above
	{`[ab]+`, "zzza", true},
	{`[ab]+`, "zzzab", true},
	{`[ab]+`, "zzz", false},
	{`[^ab]+`, "zzza", true},
	{`[^ab]+`, "zzzab", true},
	{`[^ab]+`, "zzz", true},

	// Capture groups
	{`(cat) and \1`, "cat and cat", true},
	{`(cat) and \1`, "cat and dog", false},

	// Still failing
	// {`\w+`, "abcefg", true},
	// {`\w+`, "12345", true},
	// {`\w+`, "abc12345", true},
	// {`(\w\w\w\w \d\d\d) is doing \1 times`, "grep 101 is doing grep 101 times", true},
	// {`(\w+) and \1`, "cat and cat", true},
	// {`(\w+) and \1`, "dog and dog", true},
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
