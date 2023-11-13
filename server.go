package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/jeremyKisner/streaming-daemon/audioproducer"
)

const sampleRate = 44100

var (
	requestCounter int
	mu             sync.Mutex
)

func main() {
	r := mux.NewRouter()

	// r.HandleFunc("/")
	r.HandleFunc("/", handleStream)

	http.ListenAndServe(":8080", r)
}

func handleStream(w http.ResponseWriter, r *http.Request) {
	Producer := audioproducer.Producer{}
	Producer.StreamRandomBeeps(w, r)
}
