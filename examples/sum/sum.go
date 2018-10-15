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
		Min:     1,
	}

	pool, err := twerk.New(func(a, b int) int {
		time.Sleep(1 * time.Second)
		return a + b
	}, config)

	if err != nil {
		fmt.Println(err)
		return
	}

	go exampleWork(pool, 1, 2)
	go exampleWork(pool, 2, 9)
	time.Sleep(3 * time.Second)
	go exampleWork(pool, 22, 1)
	go exampleWork(pool, 0, 4)
	go exampleWork(pool, 1, 5)
	go exampleWork(pool, 6, 7)
	go exampleWork(pool, 8, 9)

	pool.Wait()
	pool.Stop()
}

func exampleWork(pool twerk.Twerker, a, b int) {
	resChan, err := pool.Work(a, b)

	if err != nil {
		fmt.Println(err)
		return
	}
	// we know it's an int
	result, _ := (<-resChan)[0].(int)
	fmt.Printf("%d + %d = %d\n", a, b, result)
}
