package search

func hashSearch(arr map[int]string, key int) (string, bool) {
	val, found := arr[key]
	return val, found
}
