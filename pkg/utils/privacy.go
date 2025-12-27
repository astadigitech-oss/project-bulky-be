package utils

import "strings"

// MaskName masks buyer name for privacy in public reviews
// Examples:
// "John Doe" -> "J***n D**"
// "Budi" -> "B***i"
// "Ali" -> "A*i"
func MaskName(name string) string {
	if name == "" {
		return ""
	}

	words := strings.Split(name, " ")
	masked := make([]string, len(words))

	for i, word := range words {
		if len(word) <= 2 {
			// Keep short words as is
			masked[i] = word
		} else {
			// Keep first and last char, mask middle
			runes := []rune(word)
			middle := strings.Repeat("*", len(runes)-2)
			masked[i] = string(runes[0]) + middle + string(runes[len(runes)-1])
		}
	}

	return strings.Join(masked, " ")
}
