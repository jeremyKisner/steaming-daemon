package audioproducer

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type Producer struct{}

const sampleRate = 44100
const frequency = 440

func (p *Producer) EncodeWAV(fileName string, data []int) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	e := wav.NewEncoder(f, sampleRate, 16, 1, 1)

	buf := &audio.IntBuffer{Data: data, Format: &audio.Format{SampleRate: sampleRate, NumChannels: 1}}
	if err := e.Write(buf); err != nil {
		panic(err)
	}
	if err := e.Close(); err != nil {
		panic(err)
	}
}

func (p *Producer) GenerateBeep() []int {
	out := make([]int, sampleRate*30)
	for i := range out {
		out[i] = int(32767.0 * math.Sin(float64(i)*2.0*math.Pi*frequency/sampleRate))
	}
	return out
}

func (p *Producer) GenerateSong() []int {
	// Define the frequencies for the notes C, D, E, F, G, A, B
	notes := []float64{261.63, 293.66, 329.63, 349.23, 392.00, 440.00, 493.88}

	// Repeat the sequence 9 times
	repeats := 9

	out := make([]int, sampleRate*len(notes)*repeats)
	for r := 0; r < repeats; r++ {
		for j, freq := range notes {
			for i := 0; i < sampleRate; i++ {
				out[(r*len(notes)+j)*sampleRate+i] = int(32767.0 * math.Sin(float64(i)*2.0*math.Pi*freq/sampleRate))
			}
		}
	}
	return out
}

func (p *Producer) StreamSong(w http.ResponseWriter) <-chan int {
	// Create a channel to stream the music
	stream := make(chan int)

	go func() {
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
		audioBytes, err := ioutil.ReadFile(tmpfile.Name())
		if err != nil {
			fmt.Println(err)
			return
		}

		// Write the audio data to the http.ResponseWriter
		if _, err := w.Write(audioBytes); err != nil {
			fmt.Println("Error writing response: ", err)
		}
	}()

	// Return the music stream
	return stream
}
