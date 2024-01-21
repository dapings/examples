# backup request

`backup request` 主要是为了解决长尾请求，就是说处理时间比较长的请求，它的具体场景如下：

> 客户端超时配置不合理，给了一个比较大的值，比如5s。
> 请求发到服务端后，服务端因为资源问题或者其它问题，请求被阻塞无法返回。
> 结果：5s后客户端超时收到错误，发起重试。

可以看到，传统的重试也可以解决问题，但是从客户端来说这个请求的时延会比较大。

当后端阻塞后，其实可以有两种情况：

1. 一部分后端实例资源不足
2. 全部后端实例资源不足

对于第2种情况来说，普通重试和backup request都无法解决，但是都有一些正向收益和负向收益。

对于backup request的正、负向收益：

- 正向收益：资源在过载边缘，收到重试请求后一小段时间有资源释放，可以立即处理
- 负向收益：资源严重过载，收到重试请求只能使得后端过载更严重

但是正常情况下，资源都是充足的，而且现在的线上服务大多都会有限流，熔断等措施，所以总的来说，通过对资源的消耗来保证低延时，它的正向收益远大于负向收益。

虽然backup request需求很简单，但是还是有很多细节问题有待考虑，包括模块解耦，接口设计，防御编程等。

## 拥有的能力

拥有的能力包括：

- 重试时机算法
    - 固定时长
    - 二进制退避抖动
- 重试准入
- "非法"请求错误回调

## 使用示例

```go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/backup-request/request"
	"github.com/dapings/examples/go-programing-tour-2023/backup-request/request/retrygroup"
	"github.com/dapings/examples/go-programing-tour-2023/backup-request/retry"
)

func NewServer() *httptest.Server {
	rand.NewRand(NewSource(time.Now().UnixNano()))

	var (
		cnt int32
	)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		num := atomic.AddInt32(&cnt, 1)
		base := rand.Int63() % 20
		fmt.Printf("server: req seq %d, sleep %dms\n", num, base)
		time.Sleep(time.Duration(base) * time.Millisecond)
		w.Header().Set("seq", strconv.Itoa(int(num)))
		w.WriteHeader(http.StatusOK)
	}))
}

func main() {
	// 启动一个用于测试的http server
	srv := NewServer()
	defer srv.Close()

	// 构造原始请求
	task := func() (interface{}, error) {
		req, err := http.NewRequest(http.MethodHead, srv.URL, nil)
		if err != nil {
			panic(err)
		}
		return http.DefaultClient.Do(req)
	}
	// 请求回调处理函数
	// backup request相当于并发的重试，只接收第一次的响应
	// 非第一次返回的请求才会调用回调进行处理异常或回收资源等
	taskErrCallback := func(err error) {
		fmt.Printf("rece err: %v\n", err)
	}
	// 构造retry group，为构造backup request做准备
	// 返回的ctx也是构造backup request的参数
	rg, ctx := retrygroup.NewRetryGroup(context.Background(), task, retrygroup.WithErrHandler(taskErrCallback))

	var cnt int32
	// 重试时间算法
	event := request.WithEvent(retry.Fixed(time.Millisecond * 10))
	// 重试准入算法
	access := request.WithAccess(func(_ context.Context) bool {
		if atomic.LoadInt32(&cnt) > 5 {
			return false
		}
		atomic.AddInt32(&cnt, 1)
		return true
	})
	// 构造backup request
	req, err := request.NewRequest(ctx, rg, event, access)
	if err != nil {
		panic(err)
	}

	// 发起请求，并发的重试，阻塞等待结果
	resp, err := req.Do()
	if request.IsKilled(err) {
		fmt.Println("backup request is killed")
	}
	fmt.Printf("val is %v, err %v\n", resp, err)
}
```