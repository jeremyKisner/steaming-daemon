package main

import "github.com/jeremyKisner/streaming-daemon/internal"

func main() {
	Producer := internal.NewAudioProducer()
	song := Producer.GenerateSong()
	Producer.EncodeWAV("./assets/song.wav", song)
}
