/*
Package log provides a handler that logs each request/response (time, duration, status, method, path).

Make sure to include this handler above any other handler to get accurate logs.
*/
package log

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// An Handler provides a clever gzip compressing handler.
type Handler struct {
	Next http.Handler
}

// Handle returns a Handler wrapping another http.Handler.
func Handle(h http.Handler) *Handler {
	return &Handler{h}
}

// HandleFunc returns a Handler wrapping an http.HandlerFunc.
func HandleFunc(f http.HandlerFunc) *Handler {
	return Handle(f)
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
		log.Printf("  %s   %s   %s  %s", fmtDuration(start), fmtStatus(lw.status), fmtMethod(method), fmtPath(path))
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

func fmtDuration(start time.Time) string {
	return fmt.Sprintf("%s%s%13s%s", colorResetAll, colorDim, time.Since(start), colorResetAll)
}

func fmtStatus(status int) string {
	color := colorWhite

	switch {
	case status >= 200 && status <= 299:
		color += colorBackgroundGreen
	case status >= 300 && status <= 399:
		color += colorBackgroundCyan
	case status >= 400 && status <= 499:
		color += colorBackgroundYellow
	default:
		color += colorBackgroundRed
	}

	return fmt.Sprintf("%s%s %3d %s", colorResetAll, color, status, colorResetAll)
}

func fmtMethod(method string) string {
	var color string

	switch method {
	case "GET":
		color += colorGreen
	case "POST":
		color += colorCyan
	case "PUT", "PATCH":
		color += colorBlue
	case "DELETE":
		color += colorRed
	}

	return fmt.Sprintf("%s%s%s%s", colorResetAll, color, method, colorResetAll)
}

func fmtPath(path string) string {
	return fmt.Sprintf("%s%s%s%s", colorResetAll, colorDim, path, colorResetAll)
}
