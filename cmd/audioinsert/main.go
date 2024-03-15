package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

/*
This script is available for automating file uploads from command line.

	go run cmd/audio/insert/main.go -filepath examples/song.wav -name "to be titled" -artist "me" -album "album title"
*/
var (
	soundName, artist, album, filePath string
)

func main() {
	flag.StringVar(&filePath, "filepath", "", "required: filepath of audio")
	flag.StringVar(&soundName, "name", "", "optional: name of audio")
	flag.StringVar(&artist, "artist", "", "optional: artist of audio")
	flag.StringVar(&album, "album", "", "optional: album of audio")
	flag.Parse()

	// open the audio file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// create a new HTTP POST request
	url := "http://localhost:8082/audio/insert"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("audioFile", filePath)
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file to form file:", err)
		return
	}
	writer.WriteField("name", soundName)
	writer.WriteField("artist", artist)
	writer.WriteField("album", album)
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing multipart writer:", err)
		return
	}

	// create the HTTP request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// read and print the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response:", string(responseBody))
}
