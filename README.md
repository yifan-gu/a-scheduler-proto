a-scheduler-proto
=================

#Schedule Policy.

* In the branch `aging`:

  Prefer old request (avoid starvation)

* In branch `throughput`:

  Serve the small tasks first.

#To run

```shell
$ export GOPATH=$PWD
$ cd src
$ go run main.go > result.log
```

It will generate random resource and request sequences to the scheduler.

The result will be displayed in result.log

There is a discussion on the scheduling policy 
