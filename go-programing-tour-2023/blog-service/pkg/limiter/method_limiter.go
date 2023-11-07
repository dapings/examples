package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// 对某一部分接口进行流量调控

type MethodLimiter struct {
	*Limiter
}

func NewMethodLimiter() LimiterIface {
	l := &Limiter{limiterBuckets: make(map[string]*ratelimit.Bucket)}
	return MethodLimiter{
		Limiter: l,
	}
}

func (l MethodLimiter) Key(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (l MethodLimiter) GetBucket(key string) (*ratelimit.Bucket, bool) {
	bucket, ok := l.limiterBuckets[key]
	return bucket, ok
}

func (l MethodLimiter) AddBucket(rules ...LimiterBucketRule) LimiterIface {
	// TODO implement me
	panic("implement me")
}
