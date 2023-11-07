package limiter

// 简单的限流器，主要功能是对路由进行限流

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// LimiterIface 通用接口，限流器必须的方法
type LimiterIface interface {
	Key(ctx *gin.Context)
	GetBucket(key string) (*ratelimit.Bucket, bool)
	AddBucket(rules ...LimiterBucketRule) LimiterIface
}

type Limiter struct {
	limiterBuckets map[string]*ratelimit.Bucket
}

type LimiterBucketRule struct {
	// 键值对名称
	Key string
	// 间隔多久放N个令牌
	FillInterval time.Duration
	// 令牌桶的容量
	Capacity int64
	// 每次达到间隔时间后所放的具体令牌数量
	Quantum int64
}
