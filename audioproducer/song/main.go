package main

import "github.com/jeremyKisner/streaming-daemon/audioproducer"

func main() {
	Producer := &audioproducer.Producer{}
	// Generate a song
	song := Producer.GenerateSong()
	// Encode the song into a WAV file
	Producer.EncodeWAV("./assets/song.wav", song)
}
