package scheduler

//import (
//	"container/heap"
//	//"fmt"
//	//"time"
//)

//const (
//	timeDelta    = 1000000000 * 10 // 10s
//	priorityDiff = 2
//)

type Request struct {
	id     int
	Demand int
}

func NewRequest(id, demand int) *Request {
	return &Request{
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
