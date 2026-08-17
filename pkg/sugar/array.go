package sugar

func SliceLast[T comparable](c []T) T {
	return c[len(c)-1]
}

func SliceOf[T any](val ...T) []T {
	return val
}

func InArrayContain[T comparable](item T, items []T) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}

	return false
}
