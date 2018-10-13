package main

import (
	"fmt"
	"github.com/vizualni/twerk"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {

	pool, _ := twerk.New(
		func(url string) {
			fmt.Println("Started", url)
			res, _ := http.Get(url)
			ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			fmt.Println("Downloaded")
		},
		twerk.DefaultConfig,
	)

	pool.Work("http://google.com")
	pool.Work("http://google.com")
	pool.Work("http://google.com")
	pool.Work("http://google.com")
	pool.Work("http://google.com")
	pool.Work("http://google.com")
	pool.Work("http://google.com")
	pool.Work("http://google.com")
	pool.Work("http://google.com")
	pool.Work("http://google.com")
	pool.Work("http://google.com")
	pool.Work("http://google.com")
	pool.Work("http://google.com")
	pool.Work("http://google.com")
	pool.Work("http://google.com")
	pool.Work("http://google.com")
	pool.Work("http://google.com")

	time.Sleep(50 * time.Second)
}
