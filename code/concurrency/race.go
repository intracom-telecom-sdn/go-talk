package main

import "fmt"

//race start OMIT
var counter int

func main() {
	done := make(chan bool)
	numGoroutines := 42
	counterInc := 1000

	for i := 0; i < numGoroutines; i++ {
		go incrementCounter(counterInc, done) // HL
	}

	for i := 0; i < numGoroutines; i++ { <-done } // HL

	fmt.Println("Counter =", counter)
}

func incrementCounter(N int, done chan bool) {
	for i := 0; i < N; i++ {
		counter++
	}
	done <- true // HL
}

//race end OMIT
