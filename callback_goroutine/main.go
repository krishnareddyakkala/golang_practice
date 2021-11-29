package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	go m1()
	//time.Sleep(5 * time.Second)
	fmt.Println("main done")
}

func m1() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		c1(1, 2)
	}()
	wg.Wait()
	fmt.Println("m1 done")
}

func c1(a, b int) {
	time.Sleep(2 * time.Second)
	fmt.Println(a)
	fmt.Println("callback done")
}
