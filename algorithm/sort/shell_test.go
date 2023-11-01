package sort

import (
	"testing"
)

func TestShellSort(t *testing.T) {
	arr := []int{64, 25, 12, 22, 11}
	shellSort(arr)
	t.Log(arr)
}
