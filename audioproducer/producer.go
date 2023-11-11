package audioproducer

import (
	"math"
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
