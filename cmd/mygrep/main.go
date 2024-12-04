package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2)
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok := match(pattern, string(line))

	if !ok {
		os.Exit(1)
	}
}

func match(regex, text string) bool {
	if regex[0] == '^' {
		return matchNext(regex[1:], text, []string{})
	}

	for {
		if matchNext(regex, text, []string{}) {
			return true
		}

		if text == "" {
			return false
		}

		text = text[1:]
	}
}

func matchNext(regex, text string, captures []string) bool {
	if regex == "" {
		return true
	}

	if text == "" {
		return regex[0] == '$'
	}

	if regex[0] == '(' {
		alts, nextIdx := captureAlternation(regex)

		for _, alt := range alts {
			newRegex := alt + regex[nextIdx:]
			newCaptures := copyAndAppend(captures, alt)

			if matchNext(newRegex, text, newCaptures) {
				return true
			}
		}

		return false
	}

	if regex[0] == '[' {
		charSet, negated, nextIdx := captureCharSet(regex)

		if len(regex) > nextIdx {
			if regex[nextIdx] == '+' {
				if negated {
					return !matchCharSet(charSet, text[0]) && matchNext(regex[nextIdx+1:], text[1:], captures)
				}

				return matchCharSet(charSet, text[0]) && matchNext(regex[nextIdx+1:], text[1:], captures)
			}
		}

		if negated {
			return !matchCharSet(charSet, text[0]) && matchNext(regex[nextIdx:], text[1:], captures)
		}

		return matchCharSet(charSet, text[0]) && matchNext(regex[nextIdx:], text[1:], captures)
	}

	if len(regex) >= 2 {
		if regex[0] == '\\' {
			if regex[1] == 'd' {
				return matchDigit(text[0]) && matchNext(regex[2:], text[1:], captures)
			}

			if regex[1] == 'w' {
				return matchAlphaNumeric(text[0]) && matchNext(regex[2:], text[1:], captures)
			}

			if regex[1] == '1' {
				firstRegex := captures[0]
				firstText := text[:len(firstRegex)]
				if matchNext(firstRegex, firstText, []string{}) {
					nextRegex := regex[2:]
					nextText := text[len(captures[0]):]

					return matchNext(nextRegex, nextText, captures)
				}
			}
		}

		if regex[1] == '+' {
			return matchExact(regex[0], text[0]) && matchStar(regex[0], regex[2:], text[1:], captures)
		}

		if regex[1] == '?' {
			if regex[0] != text[0] {
				return matchNext(regex[2:], text, captures)
			}

			return matchNext(regex[2:], text[1:], captures)
		}
	}

	return matchExact(regex[0], text[0]) && matchNext(regex[1:], text[1:], captures)
}

func matchExact(regexCh, textCh byte) bool {
	return regexCh == '.' || regexCh == textCh
}

func matchDigit(ch byte) bool {
	digits := []byte("0123456789")

	for _, digit := range digits {
		if ch == digit {
			return true
		}
	}

	return false
}

func matchAlpha(ch byte) bool {
	alphas := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	for _, alpha := range alphas {
		if ch == alpha {
			return true
		}
	}

	return false
}

func matchAlphaNumeric(ch byte) bool {
	return matchDigit(ch) || matchAlpha(ch) || ch == '_'
}

func matchCharSet(charSet string, ch byte) bool {
	for i := 0; i < len(charSet); i++ {
		if charSet[i] == ch {
			return true
		}
	}

	return false
}

// This function is truly magic
func matchStar(ch byte, regex, text string, captures []string) bool {
	for {
		if matchNext(regex, text, captures) {
			return true
		}

		if text == "" || (text[0] != ch && ch != '.') {
			return false
		}

		text = text[1:]
	}
}

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
