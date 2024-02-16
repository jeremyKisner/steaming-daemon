package main

import "github.com/jeremyKisner/streaming-daemon/internal/audioproducer"

func main() {
	Producer := audioproducer.NewAudioProducer()
	song := Producer.GenerateSong()
	Producer.EncodeWAV("./assets/song.wav", song)
}
