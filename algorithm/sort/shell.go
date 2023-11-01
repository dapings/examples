package sort

// 将数组按照一定的间隔分组，对每个分组进行插入排序，然后逐渐缩小间隔直到为1，最后进行一次插入排序
func shellSort(arr []int) {
	n := len(arr)
	gap := n / 2

	for gap > 0 {
		for i := gap; i < n; i++ {
			tmp := arr[i]
			j := i

			// > from min to max, < from max to min
			for j >= gap && arr[j-gap] < tmp {
				arr[j] = arr[j-gap]
				j -= gap
			}

			arr[j] = tmp
		}

		gap /= 2
	}
}
