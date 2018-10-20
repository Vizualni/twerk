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
