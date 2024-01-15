package lfu

import (
	"container/heap"

	cache "github.com/dapings/examples/go-programing-tour-2023/5-cache"
)

type entry struct {
	key    string
	value  interface{}
	weight int
	index  int
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value) + 4 + 4
}

type queue []*entry

func (q *queue) Len() int {
	return len(*q)
}

func (q *queue) Less(i, j int) bool {
	return (*q)[i].weight < (*q)[j].weight
}

func (q *queue) Swap(i, j int) {
	v := *q
	v[i], v[j] = v[j], v[i]
	v[i].index = i
	v[j].index = j
}

func (q *queue) Push(x interface{}) {
	n := len(*q)
	en := x.(*entry)
	en.index = n
	*q = append(*q, en)
}

func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	en := old[n-1]
	old[n-1] = nil // avoid memory leak
	en.index = -1  // for safety
	*q = old[:n-1]
	return en
}

// update modifies the weight and value of an entry in the queue.
func (q *queue) update(en *entry, val interface{}, weight int) {
	en.value = val
	en.weight = weight
	heap.Fix(q, en.index)
}
