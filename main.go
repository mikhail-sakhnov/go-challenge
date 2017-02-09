package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/soider/go-challenge/challenge"
	"github.com/soider/go-challenge/entities"
	"github.com/soider/go-challenge/logger"
	"net/http"
	"runtime"
	"time"
)

type numHandler struct {
	timeout time.Duration
	workers int
}

func (nh numHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx, cancel := context.WithTimeout(ctx, nh.timeout)
	ctx, log := logger.WithLogger(ctx)
	workers := nh.workers
	if workers == 0 {
		workers = runtime.GOMAXPROCS(-1)
	}
	defer cancel()
	var resp entities.Response
	sanitizer := challenge.UrlsSanitizer(req.URL.Query()["u"])
	service := challenge.NewService(
		sanitizer,
		workers,
	)
	service.Start(ctx)
	rw.Header().Add("Content-Type", "application/json")
	running := true
	for {
		select {
		case state := <-service.Tick():
			resp.Numbers = state
		case <-service.Done():
			log.Print("Successfully done")
			running = false
		case <-ctx.Done():
			log.Print("Timeout happened")
			running = false
		}
		if !running {
			break
		}
	}
	enc := json.NewEncoder(rw)
	if err := enc.Encode(resp); err != nil {
		panic(err)
	}
}

func main() {
	var address string
	var timeout string
	var workers int
	flag.StringVar(&address, "address", ":8080", "Address and port to listen on")
	flag.StringVar(&timeout, "timeout", "450ms", "Handler overall timeout")
	flag.IntVar(&workers, "workers", 0, "Workers count")
	flag.Parse()
	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Print("Using default timeout")
		timeoutDuration = time.Millisecond * 450
	}
	http.Handle("/numbers", numHandler{timeout: timeoutDuration, workers: workers})
	server := http.Server{
		Addr: address,
	}
	server.ListenAndServe()
}
