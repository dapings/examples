package sort

// 选择一个基准元素（通常选择最后一个元素），将数组分割成两个子数组，其中一个子数组的所有元素都小于基准，另一个子数组的所有元素都大于基准，然后递归地对子数组进行排序，最终得到有序数组。
func quickSort(arr []int, low, high int, orderDesc bool) {
	if low < high {
		pivotIndex := partition(arr, low, high, orderDesc)
		quickSort(arr, low, pivotIndex-1, orderDesc)
		quickSort(arr, pivotIndex+1, high, orderDesc)
	}

	if low >= high {
		return
	}
}

// 创建了大量的G，伴随栈的分配，有性能损耗；同时大量的G在调度和垃圾回收检查时，也会占用一定的时间。
func quickSort_go1(arr []int, low, high int, orderDesc bool, done chan struct{}) {
	if low < high {
		pivotIndex := partition(arr, low, high, orderDesc)
		childDone := make(chan struct{}, 2)
		go quickSort_go1(arr, low, pivotIndex-1, orderDesc, childDone)
		go quickSort_go1(arr, pivotIndex+1, high, orderDesc, childDone)
		<-childDone
		<-childDone
		done <- struct{}{}
	}

	if low >= high {
		done <- struct{}{}
		return
	}
}

// 减少G的生成，并且递归深度在3以内，如果递归尝试超过3，则使用串行的方式。
func quickSort_go2(arr []int, low, high int, orderDesc bool, done chan struct{}, depth int) {
	if low >= high {
		done <- struct{}{}
		return
	}

	depth--
	pivotIndex := partition(arr, low, high, orderDesc)

	if depth > 0 {
		childDone := make(chan struct{}, 2)
		go quickSort_go2(arr, low, pivotIndex-1, orderDesc, childDone, depth)
		go quickSort_go2(arr, pivotIndex+1, high, orderDesc, childDone, depth)
		<-childDone
		<-childDone
	} else {
		quickSort(arr, low, pivotIndex-1, orderDesc)
		quickSort(arr, pivotIndex+1, high, orderDesc)
	}

	done <- struct{}{}
}

// partition  函数用于将数组划分为两部分，并返回基准元素的索引。
// 将数据分区为左右两部分。
func partition(arr []int, low, high int, orderDesc bool) int {
	pivot := arr[high] // 将最后一个值做为分界值
	i := low - 1
	for j := low; j < high; j++ {
		// < min to max, > max to min
		if orderDesc {
			if arr[j] > pivot {
				i++
				arr[i], arr[j] = arr[j], arr[i]
			}
		} else {
			// 如果小于分界值，则放到左边
			if arr[j] < pivot {
				i++
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}

	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}
