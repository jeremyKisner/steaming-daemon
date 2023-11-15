package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeremyKisner/streaming-daemon/audioproducer"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/beepstream", handleBeepStream)
	http.ListenAndServe(":8080", r)
}

func handleBeepStream(w http.ResponseWriter, r *http.Request) {
	Producer := audioproducer.Producer{}
	Producer.StreamRandomBeeps(w, r)
}
