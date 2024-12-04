package main

import (
	"strings"
)

func captureCharSet(pattern string) (string, bool, int) {
	i := 1
	charSet := ""
	negated := false

	if pattern[1] == '^' {
		i = 2
		negated = true
	}

	for ; pattern[i] != ']'; i += 1 {
		charSet += string(pattern[i])
	}

	nextIdx := i + 1

	return charSet, negated, nextIdx
}

func captureAlternation(pattern string) ([]string, int) {
	i := 0

	for ; pattern[i] != ')'; i += 1 {
	}

	all := pattern[1:i]
	return strings.Split(all, "|"), i + 1
}

func copyAndAppend(original []string, elements ...string) []string {
	newSlice := make([]string, len(original), len(original)+len(elements))
	copy(newSlice, original)

	return append(newSlice, elements...)
}
