package sort

import (
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

func TestQuickSort(t *testing.T) {
	arr := []int{64, 25, 12, 22, 11}
	// min to max
	quickSort(arr, 0, len(arr)-1, false)
	t.Log(arr)
	// max to min
	quickSort(arr, 0, len(arr)-1, true)
	t.Log(arr)
}

func TestBenchQuickSort(t *testing.T) {
	// 生成测试数据
	rand.Seed(uint64(time.Now().UnixNano()))
	n := 10000000
	testData1, testData2, testData3 := make([]int, 0, n), make([]int, 0, n), make([]int, 0, n)
	for i := 0; i < n; i++ {
		val := rand.Intn(n * 100)
		testData1 = append(testData1, val)
		testData2 = append(testData2, val)
		testData3 = append(testData3, val)
	}

	// 串行
	start := time.Now()
	quickSort(testData1, 0, len(testData1)-1, false)
	t.Log("串行执行：", time.Since(start))

	// 并发程序
	done2 := make(chan struct{})
	start = time.Now()
	go quickSort_go1(testData2, 0, len(testData2)-1, false, done2)
	<-done2
	t.Log("完全并发执行：", time.Since(start))

	done3 := make(chan struct{})
	start = time.Now()
	go quickSort_go2(testData3, 0, len(testData3)-1, false, done3, 3)
	<-done3
	t.Log("优化(控制深度)并发执行：", time.Since(start))
}
