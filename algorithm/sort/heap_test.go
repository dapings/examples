package sort

import (
	"testing"
)

func TestHeapSort(t *testing.T) {
	arr := []int{12, 11, 13, 5, 6, 7}
	heapSort(arr)
	// max heap [13 12 11 7 6 5]
	// min heap [5 6 7 11 12 13]
	t.Log(arr)
}
