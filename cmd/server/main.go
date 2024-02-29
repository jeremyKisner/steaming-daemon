package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeremyKisner/streaming-daemon/internal/database"
	"github.com/jeremyKisner/streaming-daemon/internal/handler"
	_ "github.com/lib/pq"
)

func main() {
	db, err := database.CreateConnection()
	if err != nil {
		fmt.Println("error starting app", err)
		return
	}
	defer db.Close()
	db.CreateAudioTable()
	port := ":8080"
	r := mux.NewRouter()
	r.HandleFunc("/healthz", handler.HandleHealthz)
	r.HandleFunc("/audio/insert", handler.HandleAudioInsert(db))
	r.HandleFunc("/audio/{id}", handler.HandleAudioExtraction(db))
	fmt.Printf("server started at http://localhost%s/\n", port)
	http.ListenAndServe(port, r)
}
