package main

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

const (
	gzipEncoding  = "gzip"
	flateEncoding = "deflate"
)

type compressResponseWriter struct {
	compressor io.Writer
	w          http.ResponseWriter
}

func (cw *compressResponseWriter) WriteHeader(c int) {
	cw.w.Header().Del("Content-Length")
	cw.w.WriteHeader(c)
}

func (cw *compressResponseWriter) Header() http.Header {
	return cw.w.Header()
}

func (cw *compressResponseWriter) Write(b []byte) (int, error) {
	h := cw.w.Header()
	if h.Get("Content-Type") == "" {
		h.Set("Content-Type", http.DetectContentType(b))
	}

	h.Del("Content-Length")

	return cw.compressor.Write(b)
}

func (cw *compressResponseWriter) ReadFrom(r io.Reader) (int64, error) {
	return io.Copy(cw.compressor, r)
}

type flusher interface {
	Flush() error
}

func (cw *compressResponseWriter) Flush() {
	// Flush compressed data if compressor supports it.
	if f, ok := cw.compressor.(flusher); ok {
		f.Flush()
	}
	// Flush HTTP response.
	if f, ok := cw.w.(http.Flusher); ok {
		f.Flush()
	}
}

func Compress(next http.Handler, level int) http.Handler {
	if level < gzip.DefaultCompression || level > gzip.BestCompression {
		level = gzip.DefaultCompression
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var encoding string
		for _, curEnc := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
			curEnc = strings.TrimSpace(curEnc)
			if curEnc == gzipEncoding || curEnc == flateEncoding {
				encoding = curEnc
				break
			}
		}

		if encoding == "" {
			next.ServeHTTP(w, r)
			return
		}

		var encWriter io.WriteCloser
		if encoding == gzipEncoding {
			encWriter, _ = gzip.NewWriterLevel(w, level)
		} else if encoding == flateEncoding {
			encWriter, _ = flate.NewWriter(w, level)
		}
		defer encWriter.Close()

		w.Header().Set("Content-Encoding", encoding)

		next.ServeHTTP(&compressResponseWriter{w: w, compressor: encWriter}, r)
	})
}
