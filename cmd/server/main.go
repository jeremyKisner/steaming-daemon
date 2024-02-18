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
	connector, err := database.CreateDB()
	if err != nil {
		fmt.Println("had error creating connection", err)
		return
	}
	connector.CreateTable()
	connector.ListTables()
	defer connector.Close()
	port := ":8080"
	r := mux.NewRouter()
	r.HandleFunc("/", handler.Healthz)
	r.HandleFunc("/beepstream", handler.BeepStream)
	fmt.Printf("server started at http://localhost%s/\n", port)
	http.ListenAndServe(port, r)
}
