package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sync"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

const sampleRate = 44100

var (
	requestCounter int
	mu             sync.Mutex
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCounter++
		mu.Unlock()
		fmt.Println("Request #", requestCounter)

		w.Header().Set("Content-Type", "audio/wav")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		notes := []float64{261.63, 293.66, 329.63, 349.23, 392.00, 440.00, 493.88}

		buf := &audio.IntBuffer{Data: make([]int, 0), Format: &audio.Format{SampleRate: sampleRate, NumChannels: 1}}

		for i := 0; i < 20; i++ {
			freq := notes[rand.Intn(len(notes))]
			duration := 0.1 + rand.Float64()*0.3

			for s := 0; s < int(duration*sampleRate); s++ {
				sample := int(32767.0 * math.Sin(float64(s)*2.0*math.Pi*freq/sampleRate))
				buf.Data = append(buf.Data, sample)
			}
		}

		// Create a temporary file
		tmpfile, err := os.CreateTemp("", "audio.*.wav")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer os.Remove(tmpfile.Name()) // clean up

		// Create a wav.Encoder with the temporary file as the output
		enc := wav.NewEncoder(tmpfile, sampleRate, 16, 1, 1)

		if err := enc.Write(buf); err != nil {
			fmt.Println("Error encoding WAV: ", err)
		}

		if err := enc.Close(); err != nil {
			fmt.Println("Error closing encoder: ", err)
		}

		// Read the temporary file
		audioBytes, err := os.ReadFile(tmpfile.Name())
		if err != nil {
			fmt.Println(err)
			return
		}

		// Write the audio data to the http.ResponseWriter
		if _, err := w.Write(audioBytes); err != nil {
			fmt.Println("Error writing response: ", err)
		}
	})

	http.ListenAndServe(":8080", nil)
}
