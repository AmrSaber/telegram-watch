package utils

func Contains[T comparable](arr []T, target T) bool {
	for _, value := range arr {
		if value == target {
			return true
		}
	}

	return false
}
