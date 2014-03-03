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
	Requestheap *RequestHeap

	// signal channels
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
		ResourceChan:   make(chan *Resource, defaultQueueSize),
		RequestChan:    make(chan *Request, defaultQueueSize),
		TerminateChan:  make(chan bool),
		ScheduleResult: make(chan *Result, defaultQueueSize),
	}

	heap.Init(s.Requestheap)

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
}

// return a list of node: resource
func (s *Scheduler) schedule() {
	if s.FreeResource == 0 || s.Requestheap.Len() == 0 {
		return
	}

	req := s.Requestheap.Peek()
	if s.FreeResource < req.Demand { // cluster full
		return
	}

	// create a map for storing allocation result
	result := make(map[int]int)
	s.GetResource(req.Demand, result)

	heap.Pop(s.Requestheap)

	// sendback result
	s.ScheduleResult <- &Result{req.id, req.Demand, result}
}

func (s *Scheduler) GetResource(demand int, result map[int]int) {
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
