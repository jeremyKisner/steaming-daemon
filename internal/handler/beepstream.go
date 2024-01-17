package handler

import (
	"net/http"

	"github.com/jeremyKisner/streaming-daemon/internal"
)

func BeepStream(w http.ResponseWriter, r *http.Request) {
	internal.NewAudioProducer().StreamRandomBeeps(w)
}
