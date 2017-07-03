/*
Package log provides a handler that logs each request/response (time, duration, status, method, path).

The log formatting can either be couloured or not.

Make sure to include this handler above any other handler to get accurate performance logs.
*/
package log

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	cReset    = "\033[0m"
	cDim      = "\033[2m"
	cRed      = "\033[31m"
	cGreen    = "\033[32m"
	cBlue     = "\033[34m"
	cCyan     = "\033[36m"
	cWhite    = "\033[97m"
	cBgRed    = "\033[41m"
	cBgGreen  = "\033[42m"
	cBgYellow = "\033[43m"
	cBgCyan   = "\033[46m"
)

// A handler provides a request/response logging handler.
type handler struct {
	next    http.Handler
	options *Options
}

// Options provides the handler options.
type Options struct {
	Color bool // Colors triggers a coloured formatting compatible with Unix-based terminals.
}

// Handle returns a Handler wrapping another http.Handler.
func Handle(h http.Handler, o *Options) http.Handler {
	return &handler{h, o}
}

// HandleFunc returns a Handler wrapping an http.HandlerFunc.
func HandleFunc(f http.HandlerFunc, o *Options) http.Handler {
	return Handle(f, o)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lw := &logWriter{
		ResponseWriter: w,
	}
	// Keep originals in case the response will be altered.
	method := r.Method
	path := r.URL.Path

	defer func() {
		if lw.status == 0 {
			lw.status = http.StatusOK
		}

		if h.options == nil || !h.options.Color {
			log.Printf("%s %s ▶︎ %d @ %s", method, path, lw.status, time.Since(start))
			return
		}

		var cBgStatus string
		switch {
		case lw.status >= 200 && lw.status <= 299:
			cBgStatus += cBgGreen
		case lw.status >= 300 && lw.status <= 399:
			cBgStatus += cBgCyan
		case lw.status >= 400 && lw.status <= 499:
			cBgStatus += cBgYellow
		default:
			cBgStatus += cBgRed
		}

		var cMethod string
		switch method {
		case "GET":
			cMethod += cGreen
		case "POST":
			cMethod += cCyan
		case "PUT", "PATCH":
			cMethod += cBlue
		case "DELETE":
			cMethod += cRed
		}

		log.Printf("%s  %s%13s%s   %s%s %3d %s   %s%s%s  %s%s%s", cReset, cDim, time.Since(start), cReset, cWhite, cBgStatus, lw.status, cReset, cMethod, method, cReset, cDim, path, cReset)
	}()

	h.next.ServeHTTP(lw, r)
}

// logWriter catches the status code from WriteHeader.
type logWriter struct {
	http.ResponseWriter
	status int
}

func (lw *logWriter) WriteHeader(status int) {
	if lw.status == 0 {
		lw.status = status
	}
	lw.ResponseWriter.WriteHeader(status)
}

// CloseNotify implements the http.CloseNotifier interface.
// No channel is returned if CloseNotify is not implemented by an upstream response writer.
func (lw *logWriter) CloseNotify() <-chan bool {
	n, ok := lw.ResponseWriter.(http.CloseNotifier)
	if !ok {
		return nil
	}
	return n.CloseNotify()
}

// Flush implements the http.Flusher interface.
// Nothing is done if Flush is not implemented by an upstream response writer.
func (lw *logWriter) Flush() {
	f, ok := lw.ResponseWriter.(http.Flusher)
	if ok {
		f.Flush()
	}
}

// Hijack implements the http.Hijacker interface.
// Error http.ErrNotSupported is returned if Hijack is not implemented by an upstream response writer.
func (lw *logWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := lw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return h.Hijack()
}

// Push implements the http.Pusher interface.
// http.ErrNotSupported is returned if Push is not implemented by an upstream response writer or not supported by the client.
func (lw *logWriter) Push(target string, opts *http.PushOptions) error {
	p, ok := lw.ResponseWriter.(http.Pusher)
	if !ok {
		return http.ErrNotSupported
	}
	return p.Push(target, opts)
}
