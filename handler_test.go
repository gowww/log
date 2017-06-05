package log

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLog(t *testing.T) {
	options := []*Options{nil, {Color: true}}
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	statuses := []int{200, 300, 400, 500}

	for _, option := range options {
		for _, method := range methods {
			for _, status := range statuses {
				h := HandleFunc(func(w http.ResponseWriter, r *http.Request) {
					if status != 200 {
						w.WriteHeader(status)
					}
				}, option)
				req := httptest.NewRequest(method, "/", nil)
				w := httptest.NewRecorder()
				h.ServeHTTP(w, req)

				if w.Code != status {
					t.Errorf("status code: want %v, got %v", status, w.Code)
				}
				if w.Body.String() != "" {
					t.Errorf("body: want empty, got %q", w.Body.String())
				}
			}
		}
	}
}
