package scheduler

import (
	"container/heap"
	"container/list"
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

type NodeInfo struct {
	NodeId        int
	ResourceCount int
}

func (n *NodeInfo) increaseResource(quota int) {
	n.ResourceCount = n.ResourceCount + quota
}

func (n *NodeInfo) decreaseResource(quota int) {
	n.ResourceCount = n.ResourceCount - quota
}

func (n *NodeInfo) inFreeNodes() bool {
	return n.ResourceCount != 0
}

type Scheduler struct {
	TotalResource int // total resource count
	FreeResource  int

	NodeInfos map[int]*NodeInfo // NodeId -> NodeInfo
	FreeNodes *list.List

	// request heap
	Requestheap    *RequestHeap
	RequestAgeHeap *RequestAgeHeap

	ResourceChan   chan *Resource
	RequestChan    chan *Request
	ScheduleResult chan *Result
	TerminateChan  chan bool
}

func New() *Scheduler {
	s := &Scheduler{
		NodeInfos:      make(map[int]*NodeInfo),
		FreeNodes:      list.New(),
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
	s.TotalResource = s.TotalResource + res.resourceCount
	s.FreeResource = s.FreeResource + res.resourceCount

	if _, ok := s.NodeInfos[res.nodeId]; !ok {
		// add the new node to NodeInfos and FreeNodes
		newNode := &NodeInfo{
			NodeId: res.nodeId,
		}
		s.NodeInfos[res.nodeId] = newNode

	}

	node := s.NodeInfos[res.nodeId]

	if !node.inFreeNodes() {
		s.FreeNodes.PushBack(node.NodeId)
	}
	node.increaseResource(res.resourceCount)

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

	result := make(map[int]int)
	s.GetResource(req.Demand, result)

	// just for result analysis
	aging := false
	if req.index != 0 {
		aging = true
	}

	heap.Remove(s.Requestheap, req.index)

	// sendback result
	s.ScheduleResult <- &Result{req.id, req.Demand, result, aging}
}

func (s *Scheduler) PeekOldestRequest() **Request {
	rh := s.RequestAgeHeap

	for rh.Len() > 0 {
		agingReq := rh.Peek()
		if (*agingReq).index >= 0 { // not yet served
			return agingReq
		}
		heap.Pop(rh)
	}
	return nil
}

func (s *Scheduler) GetResource(demand int, scheduledResourceMap map[int]int) {
	// TODO: wait a while for locality, Now just RR
	for {
		e := s.FreeNodes.Front()
		nodeId := e.Value.(int)

		quota := Min(s.NodeInfos[nodeId].ResourceCount, demand)

		s.NodeInfos[nodeId].decreaseResource(quota)
		s.FreeResource = s.FreeResource - quota
		demand = demand - quota

		result[nodeId] = quota
		if demand == 0 {
			break
		}
		s.FreeNodes.Remove(e)
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
