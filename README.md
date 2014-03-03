A Simple Scheduler
=================
#Structure

I used asynchronized model:
One eventloop keeps listening for new requests and resources, and dispatch these events to there handlers

#Schedule Policy

* In the branch `aging`:

  Prefer old request (avoid starvation). 

* How to deal with starvation?
  
  I maintain two heaps for requests.

  One has the smallest request on the top, 
  the other has the oldest request on the top.

  Each time choosing a request to serve, I will test whether the oldest request is too old.
  If so, I will serve that request first, otherwise I will serve the smallest request first.

* In branch `throughput`:

  Serve the small tasks first.

#To run

```shell
$ git clone git@github.com:yifan-gu/a-scheduler-proto.git
$ cd a-scheduler-proto
$ export GOPATH=$PWD
$ git checkout [branch_name]
$ go run src/main.go > result.log
```

It will generate random resource and request sequences to the scheduler.

The result will be displayed in result.log

There is a discussion on the scheduling policy 

[discussion](https://gist.github.com/yifan-gu/9286214)
