package main

import "github.com/jeremyKisner/streaming-daemon/internal/audioproducer"

func main() {
	Producer := audioproducer.NewAudioProducer()
	beep := Producer.GenerateBeep()
	Producer.EncodeWAV("./assets/beep.wav", beep)
}
