package sort

import (
	"sort"
)

// 将元素根据一定的映射关系分配到不同的桶中，然后对每个桶进行排序，最后合并所有桶的元素
// 首先找到数组中的最大值和最小值，然后计算出需要的桶的数量。
// 接下来，将数组中的元素分配到对应的桶中，并对每个桶中的元素进行排序。
// 最后，将排序后的元素放回原数组。
func bucketSort(arr []float64) {
	// 获取数组中的最大值和最小值
	minVal, maxVal := arr[0], arr[0]
	for _, val := range arr {
		if val > maxVal {
			maxVal = val
		}
		if val < minVal {
			minVal = val
		}
	}

	// 计算桶的数量
	bucketCount := int(maxVal-minVal)/len(arr) + 1
	buckets := make([][]float64, bucketCount)

	// 将元素分配到桶中
	for _, val := range arr {
		index := int(val-minVal) / len(arr)
		buckets[index] = append(buckets[index], val)
	}

	// 对每个桶中的元素进行排序
	for _, bucket := range buckets {
		sort.Float64s(bucket)
	}

	// 将排序后的元素放回原数组
	index := 0
	for _, bucket := range buckets {
		for _, val := range bucket {
			arr[index] = val
			index++
		}
	}
}
