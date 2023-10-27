package sort

// 合并两个有序数组
func merge(arr []int, left, mid, right int) {
	leftLen := mid - left + 1
	rightLen := right - mid

	// 创建临时数组
	leftArr := make([]int, leftLen)
	rightArr := make([]int, rightLen)

	// 拷贝到临时数组
	for i := 0; i < leftLen; i++ {
		leftArr[i] = arr[left+i]
	}
	for j := 0; j < rightLen; j++ {
		rightArr[j] = arr[mid+1+j]
	}

	// 合并临时数组到原始数组
	i, j, k := 0, 0, left
	for i < leftLen && j < rightLen {
		// 从大到小 >=, 从小到大 <=
		if leftArr[i] <= rightArr[j] {
			arr[k] = leftArr[i]
			i++
		} else {
			arr[k] = rightArr[j]
			j++
		}
		k++
	}

	// 处理剩余的元素
	for i < leftLen {
		arr[k] = leftArr[i]
		i++
		k++
	}
	for j < rightLen {
		arr[k] = rightArr[j]
		j++
		k++
	}
}

// 归并排序
func mergeSort(arr []int, left, right int) {
	if left < right {
		mid := (left + right) / 2

		// 分割数组，并递归排序
		mergeSort(arr, left, mid)
		mergeSort(arr, mid+1, right)

		// 合并已排序的子数组
		merge(arr, left, mid, right)
	}
}
