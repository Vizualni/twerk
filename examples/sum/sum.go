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
