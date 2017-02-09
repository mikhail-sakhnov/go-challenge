package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/soider/go-challenge/challenge"
	"github.com/soider/go-challenge/logger"
	"net/http"
	"time"
	"github.com/soider/go-challenge/entities"
)



type numHandler struct {
	timeout time.Duration
}

func (nh numHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx, cancel := context.WithTimeout(ctx, nh.timeout)
	ctx, log := logger.WithLogger(ctx)

	defer cancel()
	var resp entities.Response
	sanitizer := challenge.UrlsSanitizer(req.URL.Query()["u"])
	service := challenge.NewService(
		sanitizer,
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
	flag.StringVar(&address, "address", ":8080", "Address and port to listen on")
	flag.StringVar(&timeout, "timeout", "450ms", "Handler overall timeout")
	flag.Parse()
	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Print("Using default timeout")
		timeoutDuration = time.Millisecond * 450
	}
	http.Handle("/numbers", numHandler{timeout: timeoutDuration})
	server := http.Server{
		Addr: address,
	}
	server.ListenAndServe()
}
