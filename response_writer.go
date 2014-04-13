package bourbon

import "net/http"

type responseWriter struct {
	rw      http.ResponseWriter
	written bool
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.written = true
	return rw.rw.Write(b)
}

func (rw *responseWriter) WriteHeader(i int) {
	rw.rw.WriteHeader(i)
}

func (rw *responseWriter) Header() http.Header {
	return rw.rw.Header()
}

func (rw *responseWriter) Written() bool {
	return rw.written
}
