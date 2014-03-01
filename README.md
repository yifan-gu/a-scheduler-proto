a-scheduler-proto
=================

Just a proto

#Schedule Policy.

In branch aging:
1. Prefer old request (avoid starvation)

In branch throughput
2. Serve the small tasks first.

#To run

```shell
$ cd src
$ go run main.go
```

It will generate random resource and request sequences to the scheduler.
