package fileserve

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func AccessLog() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			le := NewLogEntry(r)
			lw := newLogResponseWriter(w, r.ProtoMajor)
			defer func() {
				le.Write(lw.Status(), lw.Header())
			}()
			next.ServeHTTP(lw, r)
		})
	}
}

type LogEntry struct {
	logger  *log.Logger
	request *http.Request
	buf     *bytes.Buffer
}

func NewLogEntry(r *http.Request) *LogEntry {
	e := &LogEntry{
		logger:  log.New(os.Stdout, "", log.LstdFlags),
		request: r,
		buf:     &bytes.Buffer{},
	}
	fmt.Fprintf(e.buf, "\"%s ", r.Method)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	fmt.Fprintf(e.buf, "%s://%s%s %s\" ", scheme, r.Host, r.RequestURI, r.Proto)

	e.buf.WriteString("from ")
	e.buf.WriteString(r.RemoteAddr)
	e.buf.WriteString(" - ")

	return e
}

func (e *LogEntry) Write(status int, header http.Header) {
	fmt.Fprintf(e.buf, "%03d", status)
	e.logger.Println(e.buf.String())
}

type logResponseWriter struct {
	http.ResponseWriter
	wroteHeader bool
	code        int
}

func newLogResponseWriter(w http.ResponseWriter, protoMajor int) *logResponseWriter {
	return &logResponseWriter{
		ResponseWriter: w,
	}
}

func (l *logResponseWriter) WriteHeader(code int) {
	if !l.wroteHeader {
		l.code = code
		l.wroteHeader = true
		l.ResponseWriter.WriteHeader(code)
	}
}

func (l *logResponseWriter) Write(buf []byte) (int, error) {
	l.WriteHeader(http.StatusOK)
	return l.ResponseWriter.Write(buf)
}

func (l *logResponseWriter) Status() int {
	return l.code
}

func (l *logResponseWriter) Flush() {
	l.wroteHeader = true
	l.ResponseWriter.(http.Flusher).Flush()
}

func (l *logResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return l.ResponseWriter.(http.Hijacker).Hijack()
}

func (l *logResponseWriter) Pusher() (pusher http.Pusher) {
	if pusher, ok := l.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}
