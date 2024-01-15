package main

import "github.com/jeremyKisner/streaming-daemon/internal"

func main() {
	Producer := &internal.AudioProducer{}
	// Generate a beep sound
	beep := Producer.GenerateBeep()
	// Encode the beep sound into a WAV file
	Producer.EncodeWAV("./assets/beep.wav", beep)
}
