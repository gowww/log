package log

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLog(t *testing.T) {
	options := []*Options{{Color: false}, {Color: true}}
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	statuses := []int{200, 300, 400, 500}

	for _, option := range options {
		for _, method := range methods {
			for _, status := range statuses {
				ts := httptest.NewServer(HandleFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(status)
					w.Write(nil)
				}, option))
				defer ts.Close()

				req, err := http.NewRequest(method, ts.URL, nil)
				if err != nil {
					t.Fatal(err)
				}

				if _, err := http.DefaultClient.Do(req); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
}
