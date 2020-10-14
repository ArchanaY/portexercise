package client

import (
	"context"
	"io"
	"net/http"
)

type Rester interface {
	Upsert(ctx context.Context, r io.ReadCloser) error
	Fetch(ctx context.Context, key string) (io.ReadCloser, error)
}

type Framework struct {
	pc Rester
}

func NewFramework(pc Rester) *Framework {
	return &Framework{pc}
}

// RespondToHeartBeat ...
func (rf Framework) RespondToHeartBeat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", "inline")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (rf Framework) PostHandler(w http.ResponseWriter, r *http.Request) {
}

func (rf Framework) GetHanlder(w http.ResponseWriter, r *http.Request) {
}
