package main

import (
	"fmt"
	"math/rand"
	"scheduler"
	"time"
)

func main() {
	closeChan := make(chan bool)
	sched := scheduler.New()
	sched.GoStart()
	go generateResource(sched.ResourceChan)
	go generateRequest(sched.RequestChan)

	go func() {
		<-time.After(time.Second * 10)
		close(closeChan)
	}()

	for {
		select {
		case m := <-sched.ScheduleResult:
			printResult(m)
		case <-closeChan:
			sched.Terminate()
			return
		}
	}
}

func printResult(res *scheduler.Result) {
	fmt.Println("For request[", res.Id, "]", "demand:", res.Demand)
	for k := range res.Alloc {
		fmt.Println("node[", k, "] resource[", res.Alloc[k], "] ")
	}
	fmt.Println("\n")
}

// TODO: goroutine leak
func generateResource(resChan chan *scheduler.Resource) {
	s := rand.NewSource(0)
	r := rand.New(s)
	for {
		number := r.Intn(100)
		res := scheduler.NewResource(number, number)
		resChan <- res
		<-time.After(time.Millisecond * 500)
	}
}

func generateRequest(reqChan chan *scheduler.Request) {
	s := rand.NewSource(0)
	r := rand.New(s)
	i := 0
	for {
		req := scheduler.NewRequest(i, r.Intn(10))
		reqChan <- req
		<-time.After(time.Millisecond * 500)
		i++
	}
}
