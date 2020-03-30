package httpmonitoring

import "net/http"

type ResponseWriterSink struct {
	http.ResponseWriter
	StatusCode int
}

func (w *ResponseWriterSink) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
