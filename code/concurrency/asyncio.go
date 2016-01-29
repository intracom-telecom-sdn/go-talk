package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

//response start OMIT
func serverSleep(w http.ResponseWriter, r *http.Request) {
	sleepTime := rand.Int() % 5
	time.Sleep(time.Duration(sleepTime) * time.Second)
	fmt.Fprintf(w, "%d", sleepTime)
}

//response end OMIT

func serverStop(w http.ResponseWriter, r *http.Request) {
	fmt.Println("server exiting...")
	os.Exit(0)
}

//main start OMIT
func main() {
	//server start OMIT
 	go func() { 
		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/", serverSleep) 
		router.HandleFunc("/stop", serverStop)
		http.ListenAndServe(":4444", router)
	}()
	//server end OMIT
	fmt.Println("Waiting for server to start...") // Don't do this
	time.Sleep(time.Duration(5) * time.Second)

	//requests start OMIT
	numRequests := 42
	done := make(chan bool) // HL

	for i := 0; i < numRequests; i++ {
		go func() { // HL
			resp, _ := http.Get("http://127.0.0.1:4444")
			respBody, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("Response came after", string(respBody), "seconds")
			done <- true // HL
		}() // HL
	}
	for i := 0; i < numRequests; i++ { <-done } // HL
	//requests end OMIT
	http.Get("http://127.0.0.1:4444/stop")
}

//main end OMIT
