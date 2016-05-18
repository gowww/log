package log

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testableContent struct {
	needsGzip bool
	body      []byte
}

func TestLog(t *testing.T) {
	options := []*Options{{Color: false}, {Color: true}}
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	statuses := []int{200, 300, 400, 500}

	for _, option := range options {
		for _, method := range methods {
			for _, status := range statuses {
				ts := httptest.NewServer(HandleFunc(option, func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(status)
					w.WriteHeader(status + 1)
					w.Write(nil)
				}))
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
