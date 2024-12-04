package main

import (
	"fmt"
	"io"
	"os"
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
	// if pattern[0] == '^' {
	// 	return matchNext(pattern[1:], line, []string{})
	// }

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
		return false
	}

	if regex[0] == '\\' {
		if regex[1] == 'd' {
			return matchDigit(text[0]) && matchNext(regex[2:], text[1:], captures)
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

// func matchNext(pattern, line string, captures []string) bool {
//		if pattern == "" {
//			return true
//		}
//
//		if pattern[0] == '$' && line == "" {
//			return true
//		}
//
//		if line == "" {
//			return false
//		}
//
//		if len(pattern) >= 2 {
//			if pattern[0] == '(' {
//				alts, nextIdx := captureAlternation(pattern)
//				for _, alt := range alts {
//					if matchNext(alt+pattern[nextIdx:], line, copyAndAppend(captures, alt)) {
//						return true
//					}
//				}
//
//				return false
//			}
//
//			if pattern[0] == '[' {
//				charSet, negated, nextIdx := captureCharSet(pattern)
//
//				if negated {
//					return !matchCharSet(charSet, line[0]) && matchNext(pattern[nextIdx:], line[1:], captures)
//				}
//
//				for _, ch := range charSet {
//					if matchNext(string(ch)+pattern[nextIdx:], line[1:], captures) {
//						return true
//					}
//				}
//
//				return false
//			}
//
//			if pattern[0] == '\\' {
//				if pattern[1] == 'd' {
//					return matchDigit(line[0]) && matchNext(pattern[2:], line[1:], captures)
//				}
//
//				if pattern[1] == 'w' {
//					return matchAlphaNumeric(line[0]) && matchNext(pattern[2:], line[1:], captures)
//				}
//
//				if pattern[1] == '1' {
//					return matchNext(captures[0]+pattern[2:], line, captures[1:])
//				}
//			}
//
//			if pattern[1] == '+' {
//				return matchLiteral(pattern[0], line[0]) && matchStar(pattern[0], pattern[2:], line[1:], captures)
//			}
//
//			if pattern[1] == '?' {
//				if pattern[0] != line[0] {
//					return matchNext(pattern[2:], line, captures)
//				}
//
//				return matchNext(pattern[2:], line[1:], captures)
//			}
//		}
//
//		if line != "" {
//			return matchLiteral(pattern[0], line[0]) && matchNext(pattern[1:], line[1:], captures)
//		}
//
//		return false
//	}
//
//	func matchLiteral(pattern, line byte) bool {
//		return pattern == line || pattern == '.'
//	}
//
// // This function is truly magic
//
//	func matchStar(ch byte, pattern, line string, captures []string) bool {
//		for {
//			if matchNext(pattern, line, captures) {
//				return true
//			}
//
//			if line == "" || (line[0] != ch && ch != '.') {
//				return false
//			}
//
//			line = line[1:]
//		}
// }

//
// func matchCharSet(charSet string, ch byte) bool {
// 	for i := 0; i < len(charSet); i++ {
// 		if charSet[i] == ch {
// 			return true
// 		}
// 	}
//
// 	return false
// }
//
// func captureCharSet(pattern string) (string, bool, int) {
// 	i := 1
// 	charSet := ""
// 	negated := false
//
// 	if pattern[1] == '^' {
// 		i = 2
// 		negated = true
// 	}
//
// 	for ; pattern[i] != ']'; i += 1 {
// 		charSet += string(pattern[i])
// 	}
//
// 	nextIdx := i + 1
//
// 	return charSet, negated, nextIdx
// }
//
// func captureAlternation(pattern string) ([]string, int) {
// 	i := 0
//
// 	for ; pattern[i] != ')'; i += 1 {
// 	}
//
// 	all := pattern[1:i]
// 	return strings.Split(all, "|"), i + 1
// }
//
// func copyAndAppend(original []string, elements ...string) []string {
// 	newSlice := make([]string, len(original), len(original)+len(elements))
// 	copy(newSlice, original)
// 	return append(newSlice, elements...)
// }
