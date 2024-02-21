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

func main() {
	var (
		soundName, artist, album, filePath string
	)
	flag.StringVar(&soundName, "name", "", "name of audio")
	flag.StringVar(&artist, "artist", "", "artist of audio")
	flag.StringVar(&album, "album", "", "album of audio")
	flag.StringVar(&filePath, "filepath", "", "filepath of audio")
	flag.Parse()

	// open the audio file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// create a new HTTP POST request
	url := "http://localhost:8080/audio/insert"
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
