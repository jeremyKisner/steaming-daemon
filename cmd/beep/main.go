package main

import "github.com/jeremyKisner/streaming-daemon/internal"

func main() {
	Producer := internal.NewAudioProducer()
	beep := Producer.GenerateBeep()
	Producer.EncodeWAV("./assets/beep.wav", beep)
}
