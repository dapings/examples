package sort

// 调整堆，使其满足堆的性质
// n 长度，i 根节点下标
func heapify(arr []int, n, i int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2

	// 找出左子节点和根节点中较大/较小的值
	// 较大 largest heap <, 较小 smallest heap >
	if left < n && arr[left] > arr[largest] {
		largest = left
	}

	// 找出右子节点和根节点中较大/较小的值
	// 较大 largest heap <, 较小 smallest heap >
	if right < n && arr[right] > arr[largest] {
		largest = right
	}
	// 如果最大值/最小值不是根节点，则进行交换，并递归调整子树
	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		heapify(arr, n, largest)
	}
}

// 通过构建最大堆和反复调整堆来实现对数组的排序
func heapSort(arr []int) {
	n := len(arr)

	// 构建最大堆，从最后一个非叶子节点开始
	for i := n/2 - 1; i >= 0; i-- {
		heapify(arr, n, i)
	}

	// 逐个将堆顶元素(最大值)移到数据未，然后重新调整堆
	for i := n - 1; i >= 0; i-- {
		arr[0], arr[i] = arr[i], arr[0]
		heapify(arr, i, 0)
	}
}
