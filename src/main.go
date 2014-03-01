package main

import (
	"fmt"
	"math/rand"
	"scheduler"
	"time"
)

const (
	totalRequest = 20
)

func main() {
	closeChan := make(chan bool)
	finishChan := make(chan bool)
	cnt := 0

	sched := scheduler.New()
	sched.GoStart()
	go generateResource(sched.ResourceChan)
	go generateRequest(sched.RequestChan)

	go func() {
		<-time.After(time.Second * 200) // timeout
		close(closeChan)
	}()

	for {
		select {
		case m := <-sched.ScheduleResult:
			printResult(m, &cnt, finishChan)
		case <-closeChan:
			fmt.Println("Timetout!")
			sched.Terminate()
			return
		case <-finishChan:
			fmt.Println("All request fulfilled!")
			sched.Terminate()
			return
		}
	}
}

func printResult(res *scheduler.Result, cnt *int, finishChan chan bool) {
	fmt.Println("\n")
	fmt.Println("Fulfill request[", res.Id, "]", "demand:", res.Demand)
	for k := range res.Alloc {
		fmt.Println("node[", k, "] resource[", res.Alloc[k], "] ")
	}
	fmt.Println("\n")

	*cnt++
	if *cnt == totalRequest {
		close(finishChan)
	}
}

// TODO: goroutine leak
func generateResource(resChan chan *scheduler.Resource) {
	s := rand.NewSource(0)
	r := rand.New(s)
	for {
		number := r.Intn(20)+1
		fmt.Printf("Get Resource for node[%d], resource count: %d\n", number, 1)

		res := scheduler.NewResource(number, number)
		resChan <- res
		<-time.After(time.Millisecond * 500)
	}
}

func generateRequest(reqChan chan *scheduler.Request) {
	s := rand.NewSource(0)
	r := rand.New(s)
	for i := 0; i < totalRequest; i++ {
		reqNum := r.Intn(10) + 1
		req := scheduler.NewRequest(i, reqNum)
		fmt.Printf("Get Request for id[%d], asking for %d resource\n", i, reqNum)
		reqChan <- req
		//<-time.After(time.Millisecond * 500)
	}
}
