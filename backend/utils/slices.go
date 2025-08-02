package utils

func DeleteElementFromSlice[T comparable](slice []T, element T) []T {
	for i, v := range slice {
		if v == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
