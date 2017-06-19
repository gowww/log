# [![gowww](https://avatars.githubusercontent.com/u/18078923?s=20)](https://github.com/gowww) log [![GoDoc](https://godoc.org/github.com/gowww/log?status.svg)](https://godoc.org/github.com/gowww/log) [![Build](https://travis-ci.org/gowww/log.svg?branch=master)](https://travis-ci.org/gowww/log) [![Coverage](https://coveralls.io/repos/github/gowww/log/badge.svg?branch=master)](https://coveralls.io/github/gowww/log?branch=master) [![Go Report](https://goreportcard.com/badge/github.com/gowww/log)](https://goreportcard.com/report/github.com/gowww/log)

Package [log](https://godoc.org/github.com/gowww/log) provides a handlers that logs each request/response (time, duration, status, method, path).  
The log formatting can either be couloured or not.

## Installing

1. Get package:

	```Shell
	go get -u github.com/gowww/log
	````

2. Import it in your code:

	```Go
	import "github.com/gowww/log"
	```

## Usage

To wrap an [http.Handler](https://golang.org/pkg/net/http/#Handler), use [Handle](https://godoc.org/github.com/gowww/log#Handle):

```Go
mux := http.NewServeMux()

mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
})

http.ListenAndServe(":8080", log.Handle(mux, nil))
````

To wrap an [http.HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc), use [HandleFunc](https://godoc.org/github.com/gowww/log#HandleFunc):

```Go
http.Handle("/", log.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}, nil))

http.ListenAndServe(":8080", nil)
```

All in all, make sure to include this handler above any other handler to get accurate performance logs.

### Colorized output

If you are on a Unix-based OS, you can get a colorized output:

```Go
log.Handle(handler, &log.Options{
	Color: true,
})
```

## Output

### Colorized

![gowww-log-color](https://user-images.githubusercontent.com/9503891/27188839-c06d7b0e-51ef-11e7-80a6-afceaab23838.png)
