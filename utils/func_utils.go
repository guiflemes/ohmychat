package utils

func Filter[T any](values []T, fn func(T) bool) []T {
	result := make([]T, 0)
	for _, value := range values {
		if fn(value) {
			result = append(result, value)
		}
	}

	return result
}

// checks if any element in the slice satisfies the given filter function.
func Any[T any](values []T, fn func(T) bool) bool {
	for _, value := range values {
		if fn(value) {
			return true
		}
	}
	return false
}
