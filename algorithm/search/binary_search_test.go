package search

import (
	"testing"
)

func TestBinarySearch(t *testing.T) {
	arr := []int{2, 3, 4, 5, 6, 10, 48}
	x := 10
	result := binarySearch(arr, x)
	if result == -1 {
		t.Logf("%d不存在", x)
	}
	except := 6
	got := result + 1
	if except != got {
		t.Fatalf("expect %d, but got %d", except, got)
	}
}
