package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
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

	CSVFile = "concurrency_timeseries.csv" // 输出的时间序列CSV
	Verbose = true                         // ★ 是否输出每个用户的行为日志（带用户ID）

	// 精确目标数（不再是概率，而是保证达到的人数）

	TargetRegister = 100
	TargetLogin    = 50
	TargetRecharge = 30
)

// ========== 全局统计 ==========
var stats struct {
	arrivalsThisSec int64 // 每秒到达计数（监控协程采样时读取并重置）

	installing     int64 // 当前正在安装的并发数（实时）
	installSuccess int64
	installFail    int64

	register int64 // 注册
	login    int64 // 登录
	recharge int64 // 充值

	// ★ 用户ID列表字段（带 mu 锁保护）
	mu              sync.Mutex
	registeredUsers []int // 注册成功的用户ID
	loggedInUsers   []int // 登录成功的用户ID
	rechargedUsers  []int // 充值成功的用户ID
}

// appendUserID 线程安全地追加用户ID到指定切片
func appendUserID(list *[]int, userID int) {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	*list = append(*list, userID)
}

// getSortedUsers 返回指定用户ID切片的排序副本（线程安全）
func getSortedUsers(list *[]int) []int {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	cp := make([]int, len(*list))
	copy(cp, *list)
	sort.Ints(cp)
	return cp
}

// Snapshot 每秒快照，用于写入CSV
type Snapshot struct {
	TimeSec        float64
	ArrivalsPerSec int64 // 每秒新增到达数

	Installing     int64
	InstallSuccess int64
	InstallFail    int64

	Register int64 // 注册
	Login    int64 // 登录
	Recharge int64 // 充值
}

var (
	snapshots  []Snapshot
	snapshotMu sync.Mutex

	// 预分配身份集合：启动前确定哪些用户会注册/登录/充值，运行时只读
	// 三个 map[int]struct{} 存储身份集合，运行时只读、无需加锁
	registerSet map[int]struct{}
	loginSet    map[int]struct{}
	rechargeSet map[int]struct{}

	// 日志输出锁（避免多goroutine交叉打印）
	logMu sync.Mutex
)

// logUser 输出带用户ID的行为日志（线程安全）
func logUser(userID int, format string, args ...interface{}) {
	if !Verbose {
		return
	}
	logMu.Lock()
	defer logMu.Unlock()
	fmt.Printf("[用户%4d] ", userID)
	fmt.Printf(format+"\n", args...)
}

// randomDuration 返回 [min, max) 之间的随机时长
func randomDuration(min, max time.Duration) time.Duration {
	if max <= min {
		return min
	}

	return min + time.Duration(rand.Int63n(int64(max-min)))
}

// generatePoissonArrivals 生成服从泊松到达分布的 n 个到达时间
//
// 原理：
//  1. 泊松过程的到达间隔服从指数分布：interval = -ln(U) / λ
//  2. λ = n / windowSeconds（平均到达率）
//  3. 累加间隔得到原始到达时间
//  4. 归一化缩放到 [0, window]，保证 n 个用户全部在窗口内到达
//     （缩放保持了间隔的相对分布形态，即前疏后密的泊松特征）
func generatePoissonArrivals(n int, window time.Duration) []time.Duration {
	arrivals := make([]time.Duration, n)
	lambda := float64(n) / window.Seconds() // 平均到达率（人/秒）

	var cumulative float64
	for i := 0; i < n; i++ {
		u := rand.Float64()
		if u < 1e-9 {
			u = 1e-9 // 避免 ln(0)
		}
		interval := -math.Log(u) / lambda // 指数分布采样
		cumulative += interval
		arrivals[i] = time.Duration(cumulative * float64(time.Second))
	}

	// 归一化缩放到窗口内，保证最后一个用户恰好在窗口末尾附近到达
	if arrivals[n-1] > 0 {
		scale := float64(window) / float64(arrivals[n-1])
		for i := range arrivals {
			arrivals[i] = time.Duration(float64(arrivals[i]) * scale)
		}
	}
	return arrivals
}

// assignTargets 精确分配注册/登录/充值用户
// 打乱所有用户ID
// 前TargetRegister个会注册
// TargetRegister中前TargetLogin个会登录
// TargetLogin中的前TargetRecharge个会充值
//
// 如 1000人（打乱顺序）
// ├── 前100人 → 标记为"会注册"
// │     ├── 前50人 → 标记为"会登录"
// │     │     └── 前30人 → 标记为"会充值"
// │     └── 后50人 → 只注册不登录
// └── 后900人 → 只安装不注册
func assignTargets(total, regCnt, loginCnt, rechargeCnt int) {
	ids := make([]int, total)
	for i := range ids {
		ids[i] = i
	}
	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})

	registerSet = make(map[int]struct{}, regCnt)
	loginSet = make(map[int]struct{}, loginCnt)
	rechargeSet = make(map[int]struct{}, rechargeCnt)

	for i := 0; i < regCnt; i++ {
		registerSet[ids[i]] = struct{}{}
		if i < loginCnt {
			loginSet[ids[i]] = struct{}{}
		}
		if i < rechargeCnt {
			rechargeSet[ids[i]] = struct{}{}
		}
	}
}

// simulateUser 模拟单个用户从到达 → 安装 → 注册 → 登录 → 充值
func simulateUser(userID int, arrival time.Duration, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()

	// 等待到达时间
	time.Sleep(arrival)

	// 到达计数：记录该用户在这一秒到达
	atomic.AddInt64(&stats.arrivalsThisSec, 1)
	//logUser(userID, "到达，开始安装")

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
		logUser(userID, "安装失败，结束")
		return
	}
	atomic.AddInt64(&stats.installSuccess, 1)
	//logUser(userID, "安装成功（耗时%v）", installTime)

	// ──── 第二步：安装成功后随机延迟，再决定是否注册 ────
	if _, ok := registerSet[userID]; !ok {
		//logUser(userID, "不在注册名单中，结束")
		return // 不在注册名单中，结束
	}
	registerDelay := randomDuration(5*time.Second, 30*time.Second)
	time.Sleep(registerDelay)
	atomic.AddInt64(&stats.register, 1)
	appendUserID(&stats.registeredUsers, userID) // 记录注册用户ID
	logUser(userID, "注册成功（安装后%v）", registerDelay)
	time.Sleep(randomDuration(500*time.Millisecond, 2*time.Second)) // 注册接口耗时

	// ──── 第三步：注册后随机延迟，再决定是否登录 ────
	if _, ok := loginSet[userID]; !ok {
		logUser(userID, "不在登录名单中，结束")
		return
	}
	loginDelay := randomDuration(5*time.Second, 20*time.Second)
	time.Sleep(loginDelay)
	atomic.AddInt64(&stats.login, 1)
	appendUserID(&stats.loggedInUsers, userID) // 记录登录用户ID
	logUser(userID, "登录成功（注册后%v）", loginDelay)
	time.Sleep(randomDuration(300*time.Millisecond, 1*time.Second))

	// ──── 第四步：登录后随机延迟，再决定是否充值 ────
	if _, ok := rechargeSet[userID]; !ok {
		logUser(userID, "不在充值名单中，结束")
		return
	}
	rechargeDelay := randomDuration(10*time.Second, 60*time.Second)
	time.Sleep(rechargeDelay)
	atomic.AddInt64(&stats.recharge, 1)
	appendUserID(&stats.rechargedUsers, userID) // 记录充值用户ID
	logUser(userID, "充值成功（登录后%v）★", rechargeDelay)
	time.Sleep(randomDuration(500*time.Millisecond, 2*time.Second))
}

// dedupTail 尾部去重：去掉末尾连续相同的采样点，只保留第一个
func dedupTail() {
	snapshotMu.Lock()
	defer snapshotMu.Unlock()

	if len(snapshots) <= 1 {
		return
	}
	last := snapshots[len(snapshots)-1]
	i := len(snapshots) - 2
	for i >= 0 {
		cur := snapshots[i]
		if cur.ArrivalsPerSec == last.ArrivalsPerSec &&
			cur.Installing == last.Installing &&
			cur.InstallSuccess == last.InstallSuccess &&
			cur.InstallFail == last.InstallFail &&
			cur.Register == last.Register &&
			cur.Login == last.Login &&
			cur.Recharge == last.Recharge {
			i--
		} else {
			break
		}
	}
	snapshots = snapshots[:i+2]
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// 启动前精确分配目标用户
	assignTargets(TotalUsers, TargetRegister, TargetLogin, TargetRecharge)

	// 生成泊松到达时间序列
	arrivals := generatePoissonArrivals(TotalUsers, ArrivalWindow)

	var wg sync.WaitGroup
	sem := make(chan struct{}, MaxConcurrent)

	start := time.Now()

	// ──── 监控协程：每秒采样一次快照，写入时间序列 ────
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// 采样前二次检查：ctx 已取消则直接退出，避免多采一个点
				select {
				case <-ctx.Done():
					return
				default:
				}

				elapsed := time.Since(start).Seconds()
				snap := Snapshot{
					TimeSec:        math.Round(elapsed*10) / 10,
					ArrivalsPerSec: atomic.SwapInt64(&stats.arrivalsThisSec, 0), // ★ 读取并重置
					Installing:     atomic.LoadInt64(&stats.installing),
					InstallSuccess: atomic.LoadInt64(&stats.installSuccess),
					InstallFail:    atomic.LoadInt64(&stats.installFail),
					Register:       atomic.LoadInt64(&stats.register),
					Login:          atomic.LoadInt64(&stats.login),
					Recharge:       atomic.LoadInt64(&stats.recharge),
				}
				snapshotMu.Lock()
				snapshots = append(snapshots, snap)
				snapshotMu.Unlock()

				// 控制台显示 "当前/目标"
				fmt.Printf("[%5.1fs] 到达:%3d/s | 安装并发:%3d | 安装成功:%4d | 注册:%3d/%d | 登录:%2d/%d | 充值:%2d/%d\n",
					snap.TimeSec, snap.ArrivalsPerSec, snap.Installing, snap.InstallSuccess,
					snap.Register, TargetRegister, snap.Login, TargetLogin, snap.Recharge, TargetRecharge)
			case <-ctx.Done():
				return
			}
		}
	}()

	// 启动1000个用户
	for i := 0; i < TotalUsers; i++ {
		wg.Add(1)
		go simulateUser(i, arrivals[i], &wg, sem)
	}

	wg.Wait()
	cancel() // 所有用户完成后立即取消监控，精确停止

	// 尾部去重（兜底：即使 ctx 竞态多采了一个点，也会被去掉）
	dedupTail()

	// ──── 写入时间序列CSV ────
	file, err := os.Create(CSVFile)
	if err != nil {
		fmt.Printf("\n创建CSV文件失败: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 表头
	_ = writer.Write([]string{
		"time_seconds", "arrivals_per_sec", "installing_concurrent", "install_success", "install_fail",
		"register", "login", "recharge",
	})

	snapshotMu.Lock()
	for _, s := range snapshots {
		_ = writer.Write([]string{
			strconv.FormatFloat(s.TimeSec, 'f', 1, 64),
			strconv.FormatInt(s.ArrivalsPerSec, 10),
			strconv.FormatInt(s.Installing, 10),
			strconv.FormatInt(s.InstallSuccess, 10),
			strconv.FormatInt(s.InstallFail, 10),
			strconv.FormatInt(s.Register, 10),
			strconv.FormatInt(s.Login, 10),
			strconv.FormatInt(s.Recharge, 10),
		})
	}
	snapshotMu.Unlock()

	elapsed := time.Since(start)

	// ========== 最终结果 ==========
	fmt.Println("\n========== 广告转化模拟结果 ==========")
	fmt.Printf("时间窗口内模拟的总用户数:       %d\n", TotalUsers)
	fmt.Printf("到达模型:       泊松分布（λ=%.2f人/秒，窗口%v）\n",
		float64(TotalUsers)/ArrivalWindow.Seconds(), ArrivalWindow)
	fmt.Printf("安装耗时范围:   %v ~ %v\n", InstallMin, InstallMax)
	fmt.Printf("最大并发安装:   %d\n", MaxConcurrent)
	fmt.Printf("安装成功:       %d  (%.1f%%)\n",
		stats.installSuccess, float64(stats.installSuccess)/TotalUsers*100)
	fmt.Printf("安装失败:       %d  (%.1f%%)\n",
		stats.installFail, float64(stats.installFail)/TotalUsers*100)
	fmt.Printf("注册用户:       %d / %d（精确目标 ✓）\n", stats.register, TargetRegister)
	fmt.Printf("登录用户:       %d / %d（精确目标 ✓）\n", stats.login, TargetLogin)
	fmt.Printf("充值用户:       %d / %d（精确目标 ✓）\n", stats.recharge, TargetRecharge)
	fmt.Println("==========================================")
	// 输出各行为的用户ID列表（排序后）
	fmt.Println("\n========== 用户ID列表 ==========")
	fmt.Printf("注册用户(%d): %v\n", len(stats.registeredUsers), getSortedUsers(&stats.registeredUsers))
	fmt.Printf("登录用户(%d): %v\n", len(stats.loggedInUsers), getSortedUsers(&stats.loggedInUsers))
	fmt.Printf("充值用户(%d): %v\n", len(stats.rechargedUsers), getSortedUsers(&stats.rechargedUsers))
	fmt.Println("========================================")

	fmt.Printf("总耗时: %v\n", elapsed)
	fmt.Printf("时间序列CSV: %s（共%d个采样点）\n", CSVFile, len(snapshots))
}
