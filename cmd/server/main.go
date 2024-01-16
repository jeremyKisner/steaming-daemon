package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeremyKisner/streaming-daemon/internal"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleHealthz)
	r.HandleFunc("/beepstream", handleBeepStream)
	fmt.Println("server started at http://localhost:8080/")
	http.ListenAndServe(":8080", r)
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	fmt.Println("health endpoint called")
	w.Write([]byte("OK"))
}

func handleBeepStream(w http.ResponseWriter, r *http.Request) {
	internal.NewAudioProducer().StreamRandomBeeps(w)
}
