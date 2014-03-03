package main

import (
	"fmt"
	"math/rand"
	"scheduler"
	"time"
)

const (
	totalRequest = 50
)

const (
	defaultResourceArrivalMillis = 10
	defaultRequstArrivalMillis   = 10
)

func main() {
	reqList := make([]*scheduler.Request, 0, totalRequest)
	resultList := make([]*scheduler.Result, 0, totalRequest)
	closeChan := make(chan bool)
	finishChan := make(chan bool)
	cnt := 0

	sched := scheduler.New()
	sched.GoStart()
	go generateResource(sched.ResourceChan)
	go generateRequest(sched.RequestChan, &reqList)

	go func() {
		<-time.After(time.Second * 200) // timeout
		close(closeChan)
	}()

	for {
		select {
		case m := <-sched.ScheduleResult:
			resultList = append(resultList, m)
			printResult(m, &cnt, finishChan)
		case <-closeChan:
			fmt.Println("Timetout!")
			printFinalResult(reqList, resultList)
			sched.Terminate()
			return
		case <-finishChan:
			fmt.Println("All requests fulfilled!")
			printFinalResult(reqList, resultList)
			sched.Terminate()
			return
		}
	}
}

func printFinalResult(reqList []*scheduler.Request, resultList []*scheduler.Result) {
	fmt.Println("\n")
	fmt.Println("Final Result:")

	fmt.Println("Request Sequence")
	for _, v := range reqList {
		fmt.Printf("Id:%d, Demand:%d\n", v.Id(), v.Demand)
	}
	fmt.Println("\n")

	fmt.Println("Result Sequence")
	for _, v := range resultList {
		fmt.Printf("Id:%d, Demand:%d\n", v.Id, v.Demand)
	}
	fmt.Println()

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
	s := rand.NewSource(1)
	r := rand.New(s)
	for {
		number := r.Intn(20) + 1
		fmt.Printf("Get Resource for node[%d], resource count: %d\n", number, 1)

		res := scheduler.NewResource(number, number)
		resChan <- res
		<-time.After(time.Millisecond * defaultResourceArrivalMillis)
	}
}

func generateRequest(reqChan chan *scheduler.Request, reqList *[]*scheduler.Request) {
	s := rand.NewSource(1)
	r := rand.New(s)
	for i := 0; i < totalRequest; i++ {
		reqNum := r.Intn(10) + 1
		req := scheduler.NewRequest(i, reqNum)
		*reqList = append(*reqList, req)
		fmt.Printf("Get Request for id[%d], asking for %d resource\n", i, reqNum)
		reqChan <- req
		<-time.After(time.Millisecond * defaultRequstArrivalMillis)
	}
}
