package scheduler

import (
	"container/heap"
	"fmt"
)

var _ = fmt.Println

const (
	defaultQueueSize = 1024
)

type Result struct {
	Id     int
	Demand int
	Alloc  map[int]int
	Aging  bool
}

type Scheduler struct {
	TotalResource int // total resource count
	FreeResource  int

	Resource map[int]int // single resource, nodeId -> resourceNumber

	NodeList []int // nodelist
	NextNode int   // index for RR

	// request heap
	Requestheap    *RequestHeap
	RequestAgeHeap *RequestAgeHeap

	ResourceChan chan *Resource
	RequestChan  chan *Request

	ScheduleResult chan *Result

	TerminateChan chan bool
}

func New() *Scheduler {
	s := &Scheduler{
		Resource:       make(map[int]int),
		Requestheap:    new(RequestHeap),
		RequestAgeHeap: new(RequestAgeHeap),
		ResourceChan:   make(chan *Resource, defaultQueueSize),
		RequestChan:    make(chan *Request, defaultQueueSize),
		TerminateChan:  make(chan bool),
		ScheduleResult: make(chan *Result, defaultQueueSize),
	}

	heap.Init(s.Requestheap)
	heap.Init(s.RequestAgeHeap)

	return s
}

// Start the event-loop
func (s *Scheduler) Start() {
	for {
		select {
		case res := <-s.ResourceChan:
			s.handleNewResource(res)
		case req := <-s.RequestChan:
			s.handleNewRequest(req)
		case <-s.TerminateChan:
			// terminate the scheduler
			return
		}
	}
}

// Nonblocking start
func (s *Scheduler) GoStart() {
	go s.Start()
}

func (s *Scheduler) handleNewResource(res *Resource) {
	// add new node
	if _, ok := s.Resource[res.nodeId]; !ok {
		s.NodeList = append(s.NodeList, res.nodeId)
	}

	// increase the resource count
	s.Resource[res.nodeId] = s.Resource[res.nodeId] + res.resourceCount
	s.TotalResource = s.TotalResource + res.resourceCount
	s.FreeResource = s.FreeResource + res.resourceCount

	s.schedule()
}

func (s *Scheduler) handleNewRequest(req *Request) {
	heap.Push(s.Requestheap, req)
	heap.Push(s.RequestAgeHeap, &req)
}

// return a list of node: resource
func (s *Scheduler) schedule() {
	if s.FreeResource == 0 || s.Requestheap.Len() == 0 {
		return
	}

	var req *Request

	agingReq := s.PeekOldestRequest()
	if agingReq != nil && (*agingReq).IsTooOld() {
		req = s.Requestheap.Get((*agingReq).index).(*Request)
	} else {
		req = s.Requestheap.Peek()
	}

	if s.FreeResource < req.Demand { // cluster full
		return
	}

	scheduledResourceMap := make(map[int]int)
	s.GetResource(req.Demand, scheduledResourceMap)

	// just for result analysis
	aging := false
	if req.index != 0 {
		aging = true
	}

	heap.Remove(s.Requestheap, req.index)

	// sendback result
	s.ScheduleResult <- &Result{req.id, req.Demand, scheduledResourceMap, aging}
}

func (s *Scheduler) PeekOldestRequest() **Request {
	rh := s.RequestAgeHeap

	for rh.Len() > 0 {
		agingReq := rh.Peek()
		if (*agingReq).index >= 0 {
			return agingReq
		}
		heap.Pop(rh)
	}
	return nil
}

func (s *Scheduler) GetResource(demand int, scheduledResourceMap map[int]int) {
	// TODO: wait a while for locality, Now just RR
	for {
		nodeId := s.NodeList[s.NextNode]

		if s.Resource[nodeId] > 0 {
			delta := Min(s.Resource[nodeId], demand)

			s.Resource[nodeId] = s.Resource[nodeId] - delta
			s.FreeResource = s.FreeResource - delta
			demand = demand - delta

			scheduledResourceMap[nodeId] = delta
			if demand == 0 {
				break
			}
		}
		s.NextNode = (s.NextNode + 1) % len(s.NodeList) // TODO: delete empty nodes
	}
}

func (s *Scheduler) Terminate() {
	close(s.TerminateChan)
}

// return the min of the two value
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
