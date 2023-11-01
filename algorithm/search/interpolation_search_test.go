package search

import (
	"testing"
)

func TestInterpolationSearch(t *testing.T) {
	arr := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	val := 30
	index := interpolationSearch(arr, val)
	if index != -1 {
		t.Logf("%d 在集合中的索引为 %d\n", val, index)
	} else {
		t.Logf("%d 不在集合中", val)
	}
}
