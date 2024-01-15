package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeremyKisner/streaming-daemon/internal"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlehealthz)
	r.HandleFunc("/beepstream", handleBeepStream)
	fmt.Println("server started at http://localhost:8080/")
	http.ListenAndServe(":8080", r)
}

func handlehealthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func handleBeepStream(w http.ResponseWriter, r *http.Request) {
	Producer := internal.AudioProducer{}
	Producer.StreamRandomBeeps(w)
}
