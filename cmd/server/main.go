package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeremyKisner/streaming-daemon/internal/handler"
	_ "github.com/lib/pq"
)

func main() {
	port := ":8080"
	r := mux.NewRouter()
	r.HandleFunc("/", handler.Healthz)
	r.HandleFunc("/beepstream", handler.BeepStream)
	fmt.Printf("server started at http://localhost%s/\n", port)
	http.ListenAndServe(port, r)
}
