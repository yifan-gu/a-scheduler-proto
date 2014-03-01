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
}

type Scheduler struct {
	TotalResource int // total resource count
	FreeResource  int

	Resource map[int]int // single resource, nodeId -> resourceNumber

	NodeList []int // nodelist
	NextNode int   // index for RR

	// request heap
	Requestheap *RequestHeap

	ResourceChan chan *Resource
	RequestChan  chan *Request

	ScheduleResult chan *Result

	TerminateChan chan bool
}

func New() *Scheduler {
	s := &Scheduler{
		Resource:       make(map[int]int),
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

	scheduledResourceMap := make(map[int]int)
	s.GetResource(req.Demand, scheduledResourceMap)

	heap.Pop(s.Requestheap)

	// sendback result
	s.ScheduleResult <- &Result{req.id, req.Demand, scheduledResourceMap}
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
