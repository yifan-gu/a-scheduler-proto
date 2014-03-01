a-scheduler-proto
=================

#Schedule Policy.

In the branch `aging`:
1. Prefer old request (avoid starvation)

In branch `throughput`
2. Serve the small tasks first.

#To run

```shell
$ export GOPATH=$PWD
$ cd src
$ go run main.go > result.log
```

It will generate random resource and request sequences to the scheduler.

The result will be displayed in result.log

There is a discussion on the scheduling policy 

[link][https://gist.github.com/yifan-gu/9286214]
