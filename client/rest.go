package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"

	"github.com/apex/log"
	"github.com/gorilla/mux"
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
	defer func() {
		if e := recover(); e != nil {
			errStr := fmt.Errorf("%v", e)
			debugStack := string(debug.Stack())
			log.Error(errStr.Error())
			fmt.Println(debugStack)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	err := rf.pc.Upsert(context.TODO(), r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (rf Framework) GetHanlder(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			errStr := fmt.Errorf("%v", e)
			debugStack := string(debug.Stack())
			log.Error(errStr.Error())
			fmt.Println(debugStack)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	params := mux.Vars(r)
	port := params["port"]

	body, err := rf.pc.Fetch(context.TODO(), port)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("Failed to fetch port: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)

	if body != nil {
		defer body.Close()
		_, err := io.Copy(w, body)
		if err != nil {
			log.Error(err.Error())
		}
	}
}
