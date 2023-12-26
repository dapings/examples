package engine

import (
	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
	"go.uber.org/zap"
)

type Scheduler interface {
	Schedule()
	Push(reqQueue ...*spider.Request)
	Pull() *spider.Request
}

type Schedule struct {
	reqChan     chan *spider.Request
	workerChan  chan *spider.Request
	priReqQueue []*spider.Request
	reqQueue    []*spider.Request
	Logger      *zap.Logger
}

func NewSchedule() *Schedule {
	s := &Schedule{}
	s.reqChan = make(chan *spider.Request)
	s.workerChan = make(chan *spider.Request)

	return s
}

func (s *Schedule) Push(reqQueue ...*spider.Request) {
	for _, req := range reqQueue {
		s.reqChan <- req
	}
}

func (s *Schedule) Pull() *spider.Request {
	return <-s.workerChan
}

func (s *Schedule) Schedule() {
	var req *spider.Request

	var ch chan *spider.Request

	for {
		if req == nil && len(s.priReqQueue) > 0 {
			req = s.priReqQueue[0]
			s.priReqQueue = s.priReqQueue[1:]
			ch = s.workerChan
		}

		if req == nil && len(s.reqQueue) > 0 {
			req = s.reqQueue[0]
			s.reqQueue = s.reqQueue[1:]
			ch = s.workerChan
		}

		select {
		case r := <-s.reqChan:
			if r.Priority > 0 {
				s.priReqQueue = append(s.priReqQueue, r)
			} else {
				s.reqQueue = append(s.reqQueue, r)
			}
		case ch <- req:
			req = nil
			ch = nil
		}
	}
}
