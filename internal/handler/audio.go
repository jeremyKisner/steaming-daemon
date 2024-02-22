package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jeremyKisner/streaming-daemon/internal/database"
	"github.com/jeremyKisner/streaming-daemon/internal/record"
)

// HandleFileUpload handles audio file uploads
func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the request
	file, handler, err := r.FormFile("audioFile")
	if err != nil {
		http.Error(w, "Error retrieving file from form", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Save the file to disk
	filePath := "/path/to/uploaded/files/" + handler.Filename
	destination, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error creating file on server", http.StatusInternalServerError)
		return
	}
	defer destination.Close()
	io.Copy(destination, file)
}

func HandleAudioInsert(db database.PostgresConnector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Println("insert audio called")
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("audioFile")
		if err != nil {
			http.Error(w, "Error retrieving file from form", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		soundName := r.FormValue("name")
		artist := r.FormValue("artist")
		album := r.FormValue("album")
		wd, err := os.Getwd()
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		filePath := wd + "/audio/" + artist + "_" + album + "_" + soundName + "_" + handler.Filename
		filePath = strings.ReplaceAll(filePath, " ", "_")
		fmt.Println(filePath)
		destination, err := os.Create(filePath)
		if err != nil {
			fmt.Println("error creating file on server", err)
			http.Error(w, "Error creating file on server", http.StatusInternalServerError)
			return
		}
		defer destination.Close()
		io.Copy(destination, file)
		fmt.Printf("Uploaded file: %s\n", handler.Filename)

		record := record.AudioRecord{
			Name:      soundName,
			Artist:    artist,
			Album:     album,
			PickupURL: filePath,
		}
		if db.InsertNewAudioRecord(record) {
			w.Write([]byte("Insert Success"))
		} else {
			w.Write([]byte("Insert Failed. Please contact Admin."))
			return
		}
	}
}

func HandleAudioExtraction(db database.PostgresConnector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strID := vars["id"]
		id, err := strconv.Atoi(strID)
		if err != nil {
			fmt.Println("error converting int", strID, err)
			http.Error(w, "Error creating file on server", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(db.ExtractAudioByInternalID(id)))
	}
}
