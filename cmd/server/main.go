package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeremyKisner/streaming-daemon/internal/handler"
	_ "github.com/lib/pq"
)

func ListTables() {
	connStr := "host=postgres user=myuser password=mypassword dbname=mydatabase sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.QueryContext(context.Background(), "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Tables:")
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(tableName)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	port := ":8080"
	ListTables()
	r := mux.NewRouter()
	r.HandleFunc("/", handler.Healthz)
	r.HandleFunc("/beepstream", handler.BeepStream)
	fmt.Printf("server started at http://localhost%s/\n", port)
	http.ListenAndServe(port, r)
}
