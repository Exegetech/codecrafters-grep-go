package main

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

//
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
//
//
//
//
