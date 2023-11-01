package search

import (
	"testing"
)

func TestHashSearch(t *testing.T) {
	arr := map[int]string{
		1: "apple",
		2: "orange",
		3: "grape",
		4: "watermelon",
		5: "banana",
		6: "leon",
	}
	key := 3
	val, found := hashSearch(arr, key)
	if found {
		t.Logf("key %d found, val: %s\n", key, val)
	} else {
		t.Logf("key %d not found\n", key)
	}
}
