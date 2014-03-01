package scheduler

import (
	//"container/heap"
	//"fmt"
	"time"
)

const (
	defaultAgeLimit = 1000000000 * 10 // 10s
)

type Request struct {
	ts     int64
	id     int
	Demand int
}

func NewRequest(id, demand int) *Request {
	return &Request{
		ts:     time.Now().UnixNano(),
		id:     id,
		Demand: demand,
	}
}

type RequestHeap []*Request

// interfaces for heap
func (rh RequestHeap) Len() int {
	return len(rh)
}

func (rh RequestHeap) Less(i, j int) bool {
	ageI := time.Now().UnixNano() - rh[i].ts
	ageJ := time.Now().UnixNano() - rh[j].ts

	if ageI > defaultAgeLimit {
		return true
	}
	if ageJ > defaultAgeLimit {
		return false
	}

	return rh[i].Demand < rh[j].Demand
}

func (rh RequestHeap) Swap(i, j int) {
	rh[i], rh[j] = rh[j], rh[i]
}

func (rh *RequestHeap) Push(u interface{}) {
	*rh = append(*rh, u.(*Request))
}

func (rh *RequestHeap) Pop() interface{} {
	old := *rh
	n := len(old)
	u := old[n-1]
	*rh = old[0 : n-1]
	return u
}

func (rh RequestHeap) Peek() interface{} {
	return rh[0]
}
