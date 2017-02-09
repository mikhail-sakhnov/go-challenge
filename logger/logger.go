package logger

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
)

type ctxKey struct {
	label string
}

var logKey = ctxKey{"logger"}

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const requestIDLength = 10

func randStringBytes() string {
	b := make([]byte, requestIDLength)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// WithLogger returns context with request-scope bounded logger in it
func WithLogger(ctx context.Context) (context.Context, *log.Logger) {
	l := log.New(os.Stdout, fmt.Sprintf("[Request-ID: %s] ", randStringBytes()), log.LstdFlags)
	return context.WithValue(ctx, logKey, l), l
}

// FromContext extracts request-scope bounded logger from context or creates new stub logger
func FromContext(ctx context.Context) *log.Logger {
	if l, ok := ctx.Value(logKey).(*log.Logger); ok {
		return l
	}
	return log.New(os.Stdout, "[No request id] ", log.LstdFlags)
}
