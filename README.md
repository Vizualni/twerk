
# Twerk, twerk, twerk!


[![Go Report Card](https://goreportcard.com/badge/github.com/Vizualni/twerk)](https://goreportcard.com/report/github.com/Vizualni/twerk)
[![CircleCI](https://circleci.com/gh/Vizualni/twerk/tree/master.svg?style=shield)](https://circleci.com/gh/Vizualni/twerk/tree/master)


[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/Vizualni/twerk/issues)

[![go doc](https://godoc.org/github.com/vizualni/twerk?status.svg)](https://godoc.org/github.com/vizualni/twerk)



## What is Twerk?

Twerk is simple worker pool for Golang that can automatically scale if needed.

### Why did I give it name Twerk?

When you say words twerk, twerk, twerk really fast it sounds like work, work, work...
And this is, after all, a twerker pool...I mean, a worker pool :bowtie:.

![Rick twerking](https://media.giphy.com/media/9homx4dDO6qu4/giphy.gif)


## Example

### Summing numbers


```go
package main

import (
	"fmt"
	"github.com/vizualni/twerk"
	"time"
)

func main() {

	config := twerk.Config{
		Refresh: 800 * time.Millisecond,
		Max:     3,
	}

	pool, err := twerk.New(func(a, b int) int {
		time.Sleep(1 * time.Second)
		return a + b
	}, config)

	if err != nil {
		fmt.Println(err)
		return
	}

	go sumWrapper(pool, 1, 2)
	go sumWrapper(pool, 2, 9)
	time.Sleep(3 * time.Second)
	go sumWrapper(pool, 22, 1)
	go sumWrapper(pool, 0, 4)
	go sumWrapper(pool, 1, 5)
	go sumWrapper(pool, 6, 7)
	go sumWrapper(pool, 8, 9)

	pool.Wait()
	fmt.Println("all done")
	pool.Stop()
	fmt.Println("workers stopped")
}

func sumWrapper(pool twerk.Twerker, a, b int) {
	resChan, _ := pool.Work(a, b)
	// we know it's an int
	result, _ := (<-resChan)[0].(int)
	fmt.Printf("%d + %d = %d\n", a, b, result)
}


```

outputs:

```
2018/10/20 13:13:48 Live: 0; Working: 0; Idle: 0, Jobs in queue: 2; Max: 3
2018/10/20 13:13:48 Starting 2 workers
2018/10/20 13:13:49 Live: 2; Working: 2; Idle: 0, Jobs in queue: 0; Max: 3
1 + 2 = 3
2 + 9 = 11
2018/10/20 13:13:50 Live: 2; Working: 0; Idle: 2, Jobs in queue: 0; Max: 3
2018/10/20 13:13:50 Stopping 2 workers
2018/10/20 13:13:50 Live: 0; Working: 0; Idle: 0, Jobs in queue: 3; Max: 3
2018/10/20 13:13:50 Starting 3 workers
2018/10/20 13:13:51 Live: 3; Working: 3; Idle: 0, Jobs in queue: 2; Max: 3
8 + 9 = 17
22 + 1 = 23
6 + 7 = 13
2018/10/20 13:13:52 Live: 3; Working: 2; Idle: 1, Jobs in queue: 0; Max: 3
2018/10/20 13:13:52 Stopping 1 workers
1 + 5 = 6
0 + 4 = 4
2018/10/20 13:13:53 Live: 2; Working: 0; Idle: 2, Jobs in queue: 0; Max: 3
2018/10/20 13:13:53 Stopping 2 workers
all done
workers stopped
```


