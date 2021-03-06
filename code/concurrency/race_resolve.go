package main

import "fmt"

//race start OMIT
func main() {
	done := make(chan bool)
	counterChan := make(chan int, 1) // HL
	numGoroutines := 42
	counterInc := 1000

	for i := 0; i < numGoroutines; i++ {
		go incrementCounter(counterInc, done, counterChan)
	}
	counterChan <- 0 // HL

	for i := 0; i < numGoroutines; i++ { <-done }

	fmt.Println("Counter =", <-counterChan) // HL
}

func incrementCounter(N int, done chan bool, cntChan chan int) {
	for i := 0; i < N; i++ {
		c := <-cntChan // HL
		c++
		cntChan <- c // HL
	}
	done <- true
}

//race end OMIT
