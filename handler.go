/*
Package log provides a handler that logs each request/response (time, duration, status, method, path).

Make sure to include this handler above any other handler to get accurate logs.
*/
package log

import (
	"log"
	"net/http"
	"time"
)

// The list of used standard Unix terminal color codes.
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

// An Handler provides a clever gzip compressing handler.
type Handler struct {
	Options *Options
	Next    http.Handler
}

// Options contains the handler options.
type Options struct {
	Color bool
}

// Handle returns a Handler wrapping another http.Handler.
func Handle(o *Options, h http.Handler) *Handler {
	return &Handler{o, h}
}

// HandleFunc returns a Handler wrapping an http.HandlerFunc.
func HandleFunc(o *Options, f http.HandlerFunc) *Handler {
	return Handle(o, f)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lw := &logWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
	// Keep originals in case the response will be altered.
	method := r.Method
	path := r.URL.Path

	defer func() {
		if h.Options == nil || !h.Options.Color {
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

	h.Next.ServeHTTP(lw, r)
}

// logWriter catches the downstream repsonse writing to get the status code.
type logWriter struct {
	http.ResponseWriter
	used   bool
	status int
}

// WriteHeader catches a downstream WriteHeader call and gets the status code.
func (lw *logWriter) WriteHeader(status int) {
	if lw.used {
		return
	}
	lw.used = true
	lw.status = status
	lw.ResponseWriter.WriteHeader(status)
}

// Write catches a downstream Write call and seals the status code.
func (lw *logWriter) Write(b []byte) (int, error) {
	lw.used = true
	return lw.ResponseWriter.Write(b)
}
