package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jeremyKisner/streaming-daemon/internal/audioproducer"
	"github.com/jeremyKisner/streaming-daemon/internal/database"
	"github.com/jeremyKisner/streaming-daemon/internal/record"
)

func HandleHealthz(w http.ResponseWriter, r *http.Request) {
	fmt.Println("health endpoint called")
	w.Write([]byte("OK"))
}

func HandleTables(db database.PostgresConnector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get tables called")
		t := db.GetTables()
		w.Write([]byte(strings.Join(t, "")))
	}
}

func BeepStream(w http.ResponseWriter, r *http.Request) {
	fmt.Println("beepstream called")
	audioproducer.NewAudioProducer().StreamRandomBeeps(w)
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
		a := db.ExtractAudioByInternalID(id)
		file, err := os.Open(a.PickupURL)
		if err != nil {
			http.Error(w, "Error opening audio file", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		contentType := getContentType(a.PickupURL)
		w.Header().Set("Content-Type", contentType)
		_, err = io.Copy(w, file)
		if err != nil {
			log.Println("Error streaming audio file:", err)
			return
		}
	}
}

func getContentType(filename string) string {
	extension := strings.ToLower(filepath.Ext(filename))
	switch extension {
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	default:
		return "application/octet-stream"
	}
}
