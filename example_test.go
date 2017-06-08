package log_test

import (
	"fmt"
	"net/http"

	"github.com/gowww/log"
)

func Example() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	})

	http.ListenAndServe(":8080", log.Handle(mux, &log.Options{
		Color: true,
	}))
}

func ExampleHandle() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	})

	http.ListenAndServe(":8080", log.Handle(mux, &log.Options{
		Color: true,
	}))
}

func ExampleHandleFunc() {
	http.Handle("/", log.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	}, &log.Options{
		Color: true,
	}))

	http.ListenAndServe(":8080", nil)
}
