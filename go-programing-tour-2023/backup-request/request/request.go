package request

import (
	"context"
	"errors"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/backup-request/retry"
)

var ErrNoAccess = errors.New("have no access to retry")

type Reqeust struct {
	retry *RetryGroup
	ctx   context.Context
	opt   option
}

func NewReqeust(ctx context.Context, retryGroup *RetryGroup, opts ...Option) (*Reqeust, error) {
	if ctx == nil {
		panic("empty context in new request")
	}

	result := &Reqeust{retry: retryGroup, ctx: ctx}
	for _, apply := range opts {
		apply(&result.opt)
	}

	if result.opt.event == nil {
		retry.Fixed(time.Second)
	}

	return result, nil
}

type option struct {
	event  func(context.Context) <-chan struct{}
	access func(context.Context) bool
}

type Option func(*option)

func WithEvent(fn func(context.Context) <-chan struct{}) Option {
	return func(o *option) {
		o.event = fn
	}
}

func WithAccess(fn func(context.Context) bool) Option {
	return func(o *option) {
		o.access = fn
	}
}

func IsKilled(err error) bool {
	return errors.Is(err, ErrNoAccess)
}

func (r *Reqeust) Do() (any, error) {
	var (
		first = make(chan struct{}, 1)
		err   = make(chan error, 1)
		event = r.opt.event(r.ctx)
	)

	go func() {
		defer func() {
			close(first)
			close(err)
		}()

		for {
			select {
			case <-r.ctx.Done():
				return
			case first <- struct{}{}:
				r.retry.Go()
			case <-event:
				for i := 0; i < 1; i++ {
					if r.opt.access == nil {
						break
					}

					if !r.opt.access(r.ctx) {
						err <- r.retry.Kill()
						return
					}
				}

				r.retry.Go()
			}
		}
	}()

	result, retryErr := r.retry.Wait()
	errGroup, ok := <-err
	if ok && errGroup == nil {
		return nil, ErrNoAccess
	}

	return result, retryErr
}
