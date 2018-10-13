
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
	"sync"
	"time"
)

func main() {

	config := twerk.Config{
		Refresh: 800 * time.Millisecond,
		Max:     3,
		Min:     1,
	}

	pool, err := twerk.New(sum, config)

	if err != nil {
		fmt.Println(err)
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(7)
	go exampleWork(pool, wg, 1, 2)
	go exampleWork(pool, wg, 2, 9)
	time.Sleep(3 * time.Second)
	go exampleWork(pool, wg, 22, 1)
	go exampleWork(pool, wg, 0, 4)
	go exampleWork(pool, wg, 5, 5)
	go exampleWork(pool, wg, 5, 5)
	go exampleWork(pool, wg, 5, 5)

	wg.Wait()
	time.Sleep(1 * time.Second)
}

func sum(a, b int) int {
	time.Sleep(1 * time.Second)

	return a + b
}

func exampleWork(pool twerk.Twerker, wg *sync.WaitGroup, a, b int) {
	resChan, err := pool.Work(a, b)

	if err != nil {
		fmt.Println(err)
		return
	}
	// we know it's an int
	result, _ := (<-resChan)[0].(int)
	fmt.Printf("%d + %d = %d\n", a, b, result)
	wg.Done()
}

```

outputs:

```
2018/10/13 19:32:03 Starting 1 workers
2018/10/13 19:32:03 Live: 1; Working: 1; Idle: 0, Jobs in queue: 1; Min: 1; Max: 3
2018/10/13 19:32:03 Starting 1 workers
2 + 9 = 11
2018/10/13 19:32:04 Live: 2; Working: 1; Idle: 1, Jobs in queue: 0; Min: 1; Max: 3
1 + 2 = 3
2018/10/13 19:32:05 Live: 2; Working: 0; Idle: 2, Jobs in queue: 0; Min: 1; Max: 3
2018/10/13 19:32:05 Stopping 1 workers
2018/10/13 19:32:06 Live: 1; Working: 1; Idle: 0, Jobs in queue: 3; Min: 1; Max: 3
2018/10/13 19:32:06 Starting 2 workers
5 + 5 = 10
2018/10/13 19:32:07 Live: 3; Working: 3; Idle: 0, Jobs in queue: 2; Min: 1; Max: 3
22 + 1 = 23
5 + 5 = 10
2018/10/13 19:32:07 Live: 3; Working: 2; Idle: 1, Jobs in queue: 0; Min: 1; Max: 3
0 + 4 = 4
5 + 5 = 10
2018/10/13 19:32:08 Live: 3; Working: 0; Idle: 3, Jobs in queue: 0; Min: 1; Max: 3
2018/10/13 19:32:08 Stopping 2 workers
```




## View count [![HitCount](http://hits.dwyl.com/Vizualni/twerk.svg)](http://hits.dwyl.com/Vizualni/twerk)


