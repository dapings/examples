package sort

import (
	"testing"
)

func TestInsertionSort(t *testing.T) {
	arr := []int{64, 25, 12, 23, 11}
	insertionSor(arr)
	t.Log(arr)
}
