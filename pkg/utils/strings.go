package utils

// Slice takes a string and returns the first n runes.
func Slice(value string, n int) string {
	var totalIteratedRunes = 0
	var substring = ""
	for _, runeValue := range value {
		if totalIteratedRunes >= n {
			break
		}
		substring += string(runeValue)
		totalIteratedRunes++
	}
	return substring
}
