package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zenazn/goji"
)

var (
	curTime time.Time
)

func main() {
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for t := range ticker.C {
			curTime = t
		}
	}()

	goji.Get("/timegen", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("%s\n", curTime))
	})

	// Listen and server on :8000 unless "PORT" environment variable is set
	goji.Serve()
}
