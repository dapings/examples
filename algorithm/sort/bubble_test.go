package sort

import (
	"fmt"
	"testing"
)

func TestBubbleSort(t *testing.T) {
	arr := []int{64, 34, 25, 12, 22, 11, 90}
	bubbleSort(arr)
	fmt.Println(arr)
}
