package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

func main() {
	var (
		soundName, artist, album string
	)
	flag.StringVar(&soundName, "name", "", "name of sound")
	flag.StringVar(&artist, "artist", "", "artist of sound")
	flag.StringVar(&album, "album", "", "album of sound")
	flag.Parse()
	requestBody := map[string]string{
		"name":   soundName,
		"artist": artist,
		"album":  album,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/audio/insert", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response:", string(responseBody))
}
