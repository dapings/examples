package sort

func bubbleSort(arr []int) {
	n := len(arr)
	// 不断比较相邻的元素并交换它们的位置来排序
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			// 从小到大
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}
