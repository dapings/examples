package sort

import (
	"testing"
)

func TestSelectionSort(t *testing.T) {
	arr := []int{64, 25, 12, 23, 11}
	selectionSort(arr)
	t.Log(arr)
}
