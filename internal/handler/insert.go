package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jeremyKisner/streaming-daemon/internal/database"
)

func InsertAudio(db database.PostgresConnector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Println("insert audio called")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		var requestBody database.InsertRequest
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		if db.InsertAudio(requestBody) {
			w.Write([]byte("Insert Success"))
		} else {
			w.Write([]byte("Insert Failed"))
		}
	}
}
