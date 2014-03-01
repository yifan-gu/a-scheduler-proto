package scheduler

import (
	"fmt"
	"time"
)

var _ = fmt.Println

const (
	DefaultTimeOutMillis = 1000 // 1000ms
)

type Request struct {
	ts     time.Time
	id     int
	Demand int
	index  int
}

func NewRequest(id, demand int) *Request {
	return &Request{
		ts:     time.Now(),
		id:     id,
		Demand: demand,
	}
}

func (r *Request) Id() int {
	return r.id
}

func (r *Request) IsTooOld() bool {
	return time.Now().After(r.ts.Add(time.Millisecond * DefaultTimeOutMillis))
}

type RequestHeap []*Request

// interfaces for heap
func (rh RequestHeap) Len() int {
	return len(rh)
}

func (rh RequestHeap) Less(i, j int) bool {
	return rh[i].Demand < rh[j].Demand
}

func (rh RequestHeap) Swap(i, j int) {
	rh[i], rh[j] = rh[j], rh[i]
	rh[i].index, rh[j].index = i, j
}

func (rh *RequestHeap) Push(u interface{}) {
	req := u.(*Request)
	req.index = rh.Len()
	*rh = append(*rh, req)
}

func (rh *RequestHeap) Pop() interface{} {
	old := *rh
	n := len(old)
	req := old[n-1]
	req.index = -1
	*rh = old[0 : n-1]
	return req
}

func (rh RequestHeap) Peek() *Request {
	return rh[0]
}

func (rh RequestHeap) Get(index int) interface{} {
	return rh[index]
}

// Another heap for maintaing aging problem
type RequestAgeHeap []**Request

// interfaces for heap
func (rh RequestAgeHeap) Len() int {
	return len(rh)
}

func (rh RequestAgeHeap) Less(i, j int) bool {
	reqI := *rh[i]
	reqJ := *rh[j]
	return reqI.ts.Before(reqJ.ts)
}

func (rh RequestAgeHeap) Swap(i, j int) {
	rh[i], rh[j] = rh[j], rh[i]
}

func (rh *RequestAgeHeap) Push(u interface{}) {
	*rh = append(*rh, u.(**Request))
}

func (rh *RequestAgeHeap) Pop() interface{} {
	old := *rh
	n := len(old)
	reqPtr := old[n-1]
	*rh = old[0 : n-1]
	return reqPtr
}

func (rh RequestAgeHeap) Peek() **Request {
	return rh[0]
}
