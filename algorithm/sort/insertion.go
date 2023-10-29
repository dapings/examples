package sort

// 将未排序的元素逐个插入到已排序部分的正确位置，直到所有元素都被插入
func insertionSor(arr []int) {
	n := len(arr)
	for i := 1; i < n; i++ {
		val := arr[i]
		j := i - 1
		// > min to max, < max to min
		for j >= 0 && arr[j] > val {
			arr[j+1] = arr[j]
			j = j - 1
		}
		arr[j+1] = val
	}
}
