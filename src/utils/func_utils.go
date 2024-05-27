package utils

import "context"

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

func Map[T any, U any](original []T, fn func(T) U) []U {
	newSlice := make([]U, 0, len(original))
	for _, item := range original {
		newSlice = append(newSlice, fn(item))
	}

	return newSlice
}

func GetUserFromContext(ctx context.Context) *User {
	user := ctx.Value(UserKey).(*User)
	return user
}

func RemoveItemByIndex[T any](slice []T, index int) []T {
	if index >= 0 && index < len(slice) {
		return append(slice[:index], slice[index+1:]...)
	}
	return slice
}
