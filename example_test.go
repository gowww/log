package log_test

import (
	"fmt"
	"net/http"

	"github.com/gowww/log"
)

func Example() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Request/response will be logged.")
	})

	http.ListenAndServe(":8080", log.Handle(mux))
}

func ExampleHandle() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Request/response will be logged.")
	})

	http.ListenAndServe(":8080", log.Handle(mux))
}

func ExampleHandleFunc() {
	http.Handle("/", log.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Request/response will be logged.")
	}))

	http.ListenAndServe(":8080", nil)
}
