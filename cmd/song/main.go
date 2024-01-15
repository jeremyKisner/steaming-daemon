package main

import "github.com/jeremyKisner/streaming-daemon/internal"

func main() {
	Producer := &internal.AudioProducer{}
	// Generate a song
	song := Producer.GenerateSong()
	// Encode the song into a WAV file
	Producer.EncodeWAV("./assets/song.wav", song)
}
