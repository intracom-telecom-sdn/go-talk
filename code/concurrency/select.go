package main

import "time"
import "fmt"

//select start OMIT
func main() {
	service1 := make(chan string, 1) // HL
	go func() {
		time.Sleep(time.Second * 1)
		service1 <- "Service 1 returned a result after 1 second."
	}()
	service2 := make(chan string, 1) // HL
	go func() {
		time.Sleep(time.Second * 42)
		service2 <- "Service 2 returned a result after 42 seconds."
	}()

	for {
		select {
		case res := <-service1: // HL
			fmt.Println(res)
		case res := <-service2: // HL
			fmt.Println(res)
		case <-time.After(time.Second * 2): // HL
			fmt.Println("Timeout waiting for services to complete after 2 seconds.")
			return
		}
	}
}

//select end OMIT
