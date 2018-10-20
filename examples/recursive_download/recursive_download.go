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
