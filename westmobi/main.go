package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TotalUsers    = 1000              // 模拟总人数
	ArrivalWindow = 120 * time.Second // 用户随机到达的时间窗口
	InstallMin    = 2 * time.Second   // 安装最短耗时
	InstallMax    = 8 * time.Second   // 安装最长耗时

	InstallRate  = 1.00 // 安装成功率
	RegisterRate = 0.10 // 安装后注册概率（目标约100人）
	LoginRate    = 0.50 // 注册后登录概率（目标约50人）
	RechargeRate = 0.60 // 登录后充值概率（目标约30人）

	MaxConcurrent = 100 // 最大同时安装并发数，真实服务器有承载上限，用带缓冲 channel 做并发闸门，超过 100 个安装请求会排队等待，避免瞬时打爆下游
)

// ========== 全局统计 ==========
var stats struct {
	installing     int64 // 当前正在安装的并发数（实时）
	installSuccess int64
	installFail    int64
	register       int64 // 注册
	login          int64 // 登录
	recharge       int64 // 充值
}

// randomDuration 返回 [min, max) 之间的随机时长
func randomDuration(min, max time.Duration) time.Duration {
	if max <= min {
		return min
	}

	return min + time.Duration(rand.Int63n(int64(max-min)))
}

// simulateUser 模拟单个用户从到达 → 安装 → 注册 → 登录 → 充值的完整路径
func simulateUser(userID int, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()

	// ──── 第0步：随机延迟到达（在时间窗口内随机开始）────
	arrivalDelay := time.Duration(rand.Int63n(int64(ArrivalWindow)))
	time.Sleep(arrivalDelay)

	// ──── 第一步：安装（受并发信号量控制）────
	sem <- struct{}{} // 获取安装槽位
	atomic.AddInt64(&stats.installing, 1)

	installTime := randomDuration(InstallMin, InstallMax)
	time.Sleep(installTime) // 模拟下载+安装耗时

	atomic.AddInt64(&stats.installing, -1)
	<-sem // 释放安装槽位

	// 判断安装是否成功
	if rand.Float64() > InstallRate {
		atomic.AddInt64(&stats.installFail, 1)
		return
	}
	atomic.AddInt64(&stats.installSuccess, 1)

	// ──── 第二步：安装成功后随机延迟，再决定是否注册 ────
	time.Sleep(randomDuration(5*time.Second, 30*time.Second))
	if rand.Float64() > RegisterRate {
		return
	}
	atomic.AddInt64(&stats.register, 1)
	time.Sleep(randomDuration(500*time.Millisecond, 2*time.Second)) // 注册接口耗时

	// ──── 第三步：注册后随机延迟，再决定是否登录 ────
	time.Sleep(randomDuration(5*time.Second, 20*time.Second))
	if rand.Float64() > LoginRate {
		return
	}
	atomic.AddInt64(&stats.login, 1)
	time.Sleep(randomDuration(300*time.Millisecond, 1*time.Second))

	// ──── 第四步：登录后随机延迟，再决定是否充值 ────
	time.Sleep(randomDuration(10*time.Second, 60*time.Second))
	if rand.Float64() > RechargeRate {
		return
	}
	atomic.AddInt64(&stats.recharge, 1)
	time.Sleep(randomDuration(500*time.Millisecond, 2*time.Second))
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	sem := make(chan struct{}, MaxConcurrent) // 并发控制信号量

	start := time.Now()

	// 立即启动1000个goroutine，每个内部先随机sleep再到达
	for i := 0; i < TotalUsers; i++ {
		wg.Add(1)
		go simulateUser(i, &wg, sem)
	}

	// ──── 实时监控：每秒打印当前安装并发数 ────
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			elapsed := time.Since(start).Seconds()
			inst := atomic.LoadInt64(&stats.installing)
			succ := atomic.LoadInt64(&stats.installSuccess)
			fail := atomic.LoadInt64(&stats.installFail)
			fmt.Printf("[%5.1fs] 安装并发: %3d | 累计成功: %4d | 累计失败: %3d\n",
				elapsed, inst, succ, fail)
		}
	}()

	wg.Wait() // 等待所有用户完成全部行为
	elapsed := time.Since(start)

	// ========== 最终结果 ==========
	println("\n========== 广告转化模拟结果 ==========")
	fmt.Printf("时间窗口内模拟的总用户数:       %d\n", TotalUsers)
	fmt.Printf("到达窗口:       %v（用户随机到达）\n", ArrivalWindow)
	fmt.Printf("安装耗时范围:   %v ~ %v\n", InstallMin, InstallMax)
	fmt.Printf("最大并发安装:   %d\n", MaxConcurrent)
	fmt.Printf("安装成功:       %d  (%.1f%%)\n",
		stats.installSuccess, float64(stats.installSuccess)/TotalUsers*100)
	fmt.Printf("安装失败:       %d  (%.1f%%)\n",
		stats.installFail, float64(stats.installFail)/TotalUsers*100)
	fmt.Printf("注册用户:       %d  (%.1f%%)  [目标约100]\n",
		stats.register, float64(stats.register)/TotalUsers*100)
	fmt.Printf("登录用户:       %d  (%.1f%%)  [目标约50]\n",
		stats.login, float64(stats.login)/TotalUsers*100)
	fmt.Printf("充值用户:       %d  (%.1f%%)  [目标约30]\n",
		stats.recharge, float64(stats.recharge)/TotalUsers*100)
	println("==========================================")
	fmt.Printf("总耗时: %v\n", elapsed)
}
