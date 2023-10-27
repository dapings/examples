package sort

import (
	"testing"
)

func TestMergeSort(t *testing.T) {
	arr := []int{64, 34, 25, 12, 22, 11, 90}
	mergeSort(arr, 0, len(arr)-1)
	// [11 12 22 25 34 64 90]
	t.Log(arr)
}
