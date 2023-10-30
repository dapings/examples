package sort

import (
	"testing"
)

func TestBucketSort(t *testing.T) {
	arr := []float64{0.42, 0.32, 0.33, 0.52, 0.37, 0.47, 0.51}
	bucketSort(arr)
	t.Log(arr)
}
