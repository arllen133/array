package polyfill

func normalize(size int, ids ...int) (start, end int) {
	if size == 0 {
		return
	}

	switch len(ids) {
	case 0:
		start = 0
		end = size
	case 1:
		start = ids[0]
		end = size
	default:
		start = ids[0]
		end = ids[1]
	}

	start = normalizeIndex(start, size)
	end = normalizeIndex(end, size)
	return
}

func normalizeIndex(index int, size int) int {
	if index < 0 {
		index += size
		if index < 0 {
			return 0
		}
	} else if index > size {
		return size
	}
	return index
}

func flattenRecursive[T any](arr []T, depth int) []T {
	var result []T
	for _, item := range arr {
		if nestedArray, ok := any(item).([]T); ok && depth > 0 {
			flattened := flattenRecursive(nestedArray, depth-1)
			result = append(result, flattened...)
		} else {
			result = append(result, item)
		}
	}
	return result
}
