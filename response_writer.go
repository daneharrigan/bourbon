package bourbon

import (
	"net/http"
	"encoding/json"
)

type responseWriter struct {
	rw      http.ResponseWriter
	written bool
	encoder    *json.Encoder
	flusher    http.Flusher
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.written = true
	return w.rw.Write(b)
}

func (w *responseWriter) WriteHeader(i int) {
	w.rw.WriteHeader(i)
}

func (w *responseWriter) Header() http.Header {
	return w.rw.Header()
}

func (w *responseWriter) Written() bool {
	return w.written
}

func (w *responseWriter) Stream(e Encodeable) {
	w.encoder.Encode(e)
	w.flusher.Flush()
}

func createResponseWriter(rw http.ResponseWriter) ResponseWriter {
	w := &responseWriter{rw: rw}
	w.encoder = json.NewEncoder(w)
	w.flusher, _ = rw.(http.Flusher)
	return w
}
