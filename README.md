A Simple Scheduler
=================
#Structure

I used asynchronized model:
One eventloop keeps listening for new requests and resources, and dispatch these events to there handlers
          

#Schedule Policy.

* In the branch `aging`:

  Prefer old request (avoid starvation)

* In branch `throughput`:

  Serve the small tasks first.
>>>>>>> 5ec416103e5e7374d7a36e7ece489282744f05c6

#To run

```shell
$ export GOPATH=$PWD
$ cd src
$ go run main.go > result.log
```

It will generate random resource and request sequences to the scheduler.

The result will be displayed in result.log

There is a discussion on the scheduling policy 

[discussion](https://gist.github.com/yifan-gu/9286214)
