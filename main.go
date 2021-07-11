package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// a handler is a function that gonna to response to a request comming to your server
func main() {

	h := http.HandlerFunc(Echo)

	log.Println("listening to localhost:8000")
	if err := http.ListenAndServe("localhost:8000", h); err != nil {
		log.Fatal(err)
	}

}

// echo just about the request that you made
func Echo(w http.ResponseWriter, r *http.Request) {

	id := rand.Intn(1000)

	fmt.Println("starting: ", id)

	time.Sleep(3 * time.Second)

	fmt.Fprintf(w, "You asked to %s %s\n", r.Method, r.URL.Path)

	fmt.Println("ending: ", id)
}
