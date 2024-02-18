package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jeremyKisner/streaming-daemon/internal/database"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Println("health endpoint called")
	w.Write([]byte("OK"))
}

func GetTables(db database.PostgresConnector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := db.GetTables()
		w.Write([]byte(strings.Join(t, "")))
	}
}
