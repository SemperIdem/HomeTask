package utils

import "sort"

// ProcessNumbers does all the things: merge new array into sorted source.
// the code doesn't look easy, but it the costs of efficiency
func ProcessNumbers(source, new []int) []int {

	sort.Ints(new)

	if len(source) == 0 {
		new = deduplicate(new)
		return new
	}

	if len(new) == 0 {
		return source
	}

	merged := make([]int, 0, len(source)+len(new))

	// Pointers for slices a and b
	i, j := 0, 0

	for i < len(source) && j < len(new) {
		if source[i] < new[j] {
			merged = appendIfNotExists(merged, source[i])
			i++
		} else if source[i] > new[j] {
			merged = appendIfNotExists(merged, new[j])
			j++
		} else {
			merged = appendIfNotExists(merged, source[i])
			i++
			j++
		}
	}

	// Append remaining elements from source
	for i < len(source) {
		merged = appendIfNotExists(merged, source[i])
		i++
	}

	// Append remaining elements from b
	for j < len(new) {
		merged = appendIfNotExists(merged, new[j])
		j++
	}

	return merged
}

func appendIfNotExists(slice []int, val int) []int {
	if len(slice) == 0 || slice[len(slice)-1] != val {
		return append(slice, val)
	}
	return slice
}

func deduplicate(s []int) []int {
	if len(s) < 2 {
		return s
	}
	e := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			continue
		}
		s[e] = s[i]
		e++
	}
	return s[:e]
}
