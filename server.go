package main

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/zenazn/goji"
)

func main() {
	goji.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "pong\n")
	})

	goji.Get("/time", func(w http.ResponseWriter, r *http.Request) {
		res, err := http.Get("http://localhost:8001/timegen")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		} else {
			defer res.Body.Close()
			cnt, _ := ioutil.ReadAll(res.Body)
			io.WriteString(w, string(cnt))
		}
	})

	// Listen and server on :8000 unless "PORT" environment variable is set
	goji.Serve()
}
