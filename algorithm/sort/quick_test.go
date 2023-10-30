package sort

import (
	"testing"
)

func TestQuickSort(t *testing.T) {
	arr := []int{64, 25, 12, 22, 11}
	quickSort(arr, 0, len(arr)-1)
	t.Log(arr)
}
