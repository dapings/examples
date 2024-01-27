package request

import (
	"context"
	"errors"
	"sync"
)

var ErrKilled = errors.New("task has done")

type RetryGroup struct {
	task func() (any, error)

	successOnce sync.Once
	done        chan struct{}
	cancel      func()

	val interface{}
	err error

	opt optionGroup
}

type optionGroup struct {
	errHandler func(error)
}

type OptionGroup func(*optionGroup)

func WithErrHandler(handler func(error)) OptionGroup {
	return func(o *optionGroup) {
		o.errHandler = handler
	}
}

func NewRetryGroup(ctx context.Context, task func() (any, error), opts ...OptionGroup) (*RetryGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	retry := &RetryGroup{cancel: cancel, task: task, done: make(chan struct{})}
	for _, apply := range opts {
		apply(&retry.opt)
	}

	return retry, ctx
}

func (g *RetryGroup) Go() {
	go func() {
		if g.task == nil {
			panic("no task in retry group")
		}

		done := false
		result, err := g.task()
		g.successOnce.Do(func() {
			g.val = result
			g.err = err

			select {
			case <-g.done:
				break
			default:
				close(g.done)

				if g.cancel != nil {
					g.cancel()
				}
			}

			done = true
		})

		if done {
			return
		}

		if g.opt.errHandler != nil {
			g.opt.errHandler(err)
		}
	}()
}

func (g *RetryGroup) Wait() (any, error) {
	if g.done == nil {
		panic("invoke Go func First")
	}

	<-g.done

	return g.val, g.err
}

func (g *RetryGroup) Kill() error {
	if g.done == nil {
		panic("invoke Go func First")
	}

	select {
	case <-g.done:
		return ErrKilled
	default:
		close(g.done)
		if g.cancel != nil {
			g.cancel()
		}
	}

	return nil
}
