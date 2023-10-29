package search

// arr 从小到大的一个有序数组, x 需要找的数据
func binarySearch(arr []int, x int) int {
	low := 0
	high := len(arr) - 1
	for low <= high {
		mid := (low + high) / 2
		if arr[mid] < x {
			low = mid + 1
		} else if arr[mid] > x {
			high = mid - 1
		} else {
			return mid
		}
	}

	// 不在数组中
	return -1
}
