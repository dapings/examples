package search

// 插值查找适用于有序且分布均匀的数据集
func interpolationSearch(arr []int, val int) int {
	low := 0
	high := len(arr) - 1

	for low <= high && val >= arr[low] && val <= arr[high] {
		pos := low + ((val-arr[low])*(high-low))/(arr[high]-arr[low])
		if arr[pos] == val {
			return pos
		} else if arr[pos] < val {
			low = pos + 1
		} else {
			high = pos - 1
		}
	}

	// 不存在时，返回 -1 索引
	return -1
}
