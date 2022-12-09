package util

func IsInArrayInt(arr []int, el int) bool {
	for _, e := range arr {
		if e == el {
			return true
		}
	}
	return false
}

func IsInArrayStr(arr []string, el string) bool {
	for _, e := range arr {
		if e == el {
			return true
		}
	}
	return false
}

func IsIntersectArrayStr(arr1, arr2 []string) bool {
	for _, e := range arr1 {
		if IsInArrayStr(arr2, e) {
			return true
		}
	}
	return false
}

func IsIntersectArrayInt(arr1, arr2 []int) bool {
	for _, e := range arr1 {
		if IsInArrayInt(arr2, e) {
			return true
		}
	}
	return false
}

// this function returns non-intersect element from first array
func FilterOutIntersectSliceStr(arr1, arr2 []string) []string {
	var result []string
	for _, e := range arr1 {
		if IsInArrayStr(arr2, e) {
			continue
		}
		result = append(result, e)
	}
	return result
}
