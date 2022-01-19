package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type GzipHandler struct {
}

func NewGzipHandler() *GzipHandler {
	return &GzipHandler{}
}

// middleware for gzip
func (g *GzipHandler) MiddleWare(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			ww := NewWrappedResponseWriter(w)
			ww.w.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(ww, r)
			defer ww.Flush()

			return
		}
		next.ServeHTTP(w, r)
	})
}

type WrappedResponseWriter struct {
	w http.ResponseWriter
	gz *gzip.Writer
}

func NewWrappedResponseWriter(w http.ResponseWriter) *WrappedResponseWriter {
	gz := gzip.NewWriter(w)
	return &WrappedResponseWriter{w: w, gz: gz}
}

// implementing interface
func (w *WrappedResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *WrappedResponseWriter) Write(d []byte) (int, error) {
	return w.gz.Write(d)
}

func (w *WrappedResponseWriter) WriteHeader(statuscode int) {
	w.w.WriteHeader(statuscode)
}

func (w *WrappedResponseWriter) Flush() {
	w.gz.Flush()
	w.gz.Close()
}
