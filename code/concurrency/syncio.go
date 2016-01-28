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
	go func() {
		//server start OMIT
		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/", serverSleep) // HL
		router.HandleFunc("/stop", serverStop)
		http.ListenAndServe(":4444", router) // HL
		//server end OMIT
	}()

	//requests start OMIT
	numRequests := 42

	for i := 0; i < numRequests; i++ {
		resp, _ := http.Get("http://127.0.0.1:4444")
		respBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response came after", string(respBody), "seconds")
	}
	//requests end OMIT

	http.Get("http://127.0.0.1:4444/stop")
}

//main end OMIT
