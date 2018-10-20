
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

#### Output

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

### Recursive crawler

```go
package main

import (
	"fmt"
	"github.com/vizualni/twerk"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

var (
	urlsRegex = regexp.MustCompile(`"(http(s?)://.+?)"`)
)

type crawler struct {
	visited sync.Map
}

func main() {

	c := &crawler{
		visited: sync.Map{},
	}

	pool, err := twerk.New(
		c.recursiveDownload,
		twerk.DefaultConfig,
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = pool.Work(pool, "http://blog.golang.com", 2)

	if err != nil {
		fmt.Println(err)
		return
	}

	pool.Wait()

	count := 0
	c.visited.Range(func(key, value interface{}) bool {
		count++
		return true
	})

	fmt.Println(count)
	pool.Stop()

}


// downloads url and extracts all the links it can find
// then sends then back to the twerker
// once the depth was hit to zero, it will exit (as this is just an example)
func (c *crawler) recursiveDownload(pool twerk.Twerker, url string, depth int) {
	if depth <= 0 {
		return
	}
	if _, visited := c.visited.Load(url); visited {
		return
	}
	c.visited.Store(url, true)

	fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		return
	}

	b, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	body := string(b)

	matches := urlsRegex.FindAllStringSubmatch(body, -1)

	for _, match := range matches {
		pool.Work(pool, match[1], depth-1)
	}

}
```

#### Output


```
2018/10/20 13:25:55 Live: 0; Working: 0; Idle: 0, Jobs in queue: 1; Max: 4
2018/10/20 13:25:55 Starting 1 workers
http://blog.golang.com
https://blog.golang.org/feed.atom
2018/10/20 13:25:56 Live: 1; Working: 1; Idle: 0, Jobs in queue: 4; Max: 4
2018/10/20 13:25:56 Starting 3 workers
https://cloud.google.com/blog/products/application-development/go-1-11-is-now-available-on-app-engine
https://www.googletagmanager.com/gtag/js?id=UA-11222381-3
https://cloud.google.com/appengine/
https://blog.golang.org/go-and-google-app-engine
https://cloud.google.com/appengine/docs/standard/go111/specifying-dependencies
https://golang.org/cmd/go/#hdr-Vendor_Directories
https://golang.org/doc/go1.11#modules
https://twitter.com/kelseyhightower/status/1035278586754813952
https://cloud.google.com/resource-manager/docs/creating-managing-projects
2018/10/20 13:25:57 Live: 4; Working: 4; Idle: 0, Jobs in queue: 4; Max: 4
https://cloud.google.com/sdk/
https://cloud.google.com/appengine/docs/standard/go111/go-differences
https://cloud.google.com/free/
https://cloud.google.com/appengine/docs/standard/go111/building-app/
https://blog.golang.org/go-cloud
https://github.com/google/go-cloud
https://en.wikipedia.org/wiki/Dependency_injection
https://cloud.google.com/open-cloud/
https://godoc.org/github.com/google/go-cloud/blob/gcsblob#OpenBucket
2018/10/20 13:25:58 Live: 4; Working: 4; Idle: 0, Jobs in queue: 4; Max: 4
https://godoc.org/github.com/google/go-cloud/blob
https://github.com/uber-go/dig
https://godoc.org/github.com/google/go-cloud/blob/s3blob
https://github.com/facebookgo/inject
https://en.wikipedia.org/wiki/Service_locator_pattern
https://google.github.io/dagger/
https://github.com/google/go-cloud/blob/master/wire/README.md
https://github.com/google/go-cloud/tree/master/samples/wire
https://github.com/google/go-cloud/issues/new
http://goo.gl/nnPfct
2018/10/20 13:25:59 Live: 4; Working: 4; Idle: 0, Jobs in queue: 4; Max: 4
https://policies.google.com/privacy
https://groups.google.com/forum/#!forum/go-cloud
https://blog.golang.org/toward-go2
https://go.googlesource.com/proposal/+/master/design/go2draft.md
https://golang.org/dl/
https://golang.org/issue/new
https://golang.org/doc/go1.11#wasm
https://golang.org/wiki/WebAssembly
https://webassembly.org/
https://golang.org/issues/new
2018/10/20 13:26:00 Live: 4; Working: 4; Idle: 0, Jobs in queue: 4; Max: 4
https://github.com/neelance
https://golang.org/doc/go1.11
https://developers.google.com/site-policies#restrictions
https://go.googlesource.com/blog/
2018/10/20 13:26:01 Live: 4; Working: 0; Idle: 4, Jobs in queue: 0; Max: 4
2018/10/20 13:26:01 Stopping 4 workers
44
```
