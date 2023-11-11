package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Open the audio file
		f, err := os.Open("." + r.URL.Path)
		if err != nil {
			http.Error(w, "File not found", 404)
			return
		}
		defer f.Close()

		// Set the Content-Type header to serve the audio file correctly
		w.Header().Set("Content-Type", "audio/wav")

		// Create a buffer to control the size of each chunk
		buf := make([]byte, 32*1024)

		fmt.Println("request received, starting streaming")
		// Stream the audio file to the client in chunks
		io.CopyBuffer(w, f, buf)
	})

	// Start the HTTP server
	http.ListenAndServe(":8080", nil)
}
