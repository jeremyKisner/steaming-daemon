package handler

import (
	"net/http"

	"github.com/jeremyKisner/streaming-daemon/internal/audioproducer"
)

func BeepStream(w http.ResponseWriter, r *http.Request) {
	audioproducer.NewAudioProducer().StreamRandomBeeps(w)
}
