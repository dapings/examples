package sort

// 选择一个基准元素（通常选择最后一个元素），将数组分割成两个子数组，其中一个子数组的所有元素都小于基准，另一个子数组的所有元素都大于基准，然后递归地对子数组进行排序，最终得到有序数组。
func quickSort(arr []int, low, high int) {
	if low < high {
		pivotIndex := partition(arr, low, high)
		quickSort(arr, low, pivotIndex-1)
		quickSort(arr, pivotIndex+1, high)
	}
}

// partition  函数用于将数组划分为两部分，并返回基准元素的索引。
func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1
	for j := low; j < high; j++ {
		// < min to max, > max to min
		if arr[j] > pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}
