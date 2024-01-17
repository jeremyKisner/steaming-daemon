package handler

import (
	"fmt"
	"net/http"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Println("health endpoint called")
	w.Write([]byte("OK"))
}
