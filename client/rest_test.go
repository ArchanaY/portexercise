package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

type mockClient struct {
	tpanic  bool
	succeed bool
	key     string
}

func (mc *mockClient) Upsert(ctx context.Context, r io.ReadCloser) error {
	if mc.succeed {
		dec := json.NewDecoder(r)
		_, _ = dec.Token()
		for dec.More() {
			key, _ := dec.Token()
			mc.key = key.(string)
			break
		}
	} else {
		return errors.New("Unknown error")
	}
	return nil
}

func (mc *mockClient) Fetch(ctx context.Context, key string) (io.ReadCloser, error) {
	return nil, nil
}

func initRESTLayers(t *testing.T, tpanic bool, succeed bool) (Framework, *mockClient) {
	t.Helper()
	ms := mockClient{succeed: succeed, tpanic: tpanic}
	rf := NewFramework(&ms)
	return *rf, &ms
}

func TestFramework_PostHandler(t *testing.T) {
	t.Run("Happy Days", func(t *testing.T) {
		rf, mock := initRESTLayers(t, false, true)
		w := httptest.NewRecorder()
		data := "{\"key2\": { \"name\": \"name2\",\"city\": \"city2\"}"
		r := httptest.NewRequest("POST", "/port", strings.NewReader(data))
		r.Header.Add("Content-Type", "application/json")
		router := mux.NewRouter()
		router.HandleFunc("/port", rf.PostHandler).Methods("POST")
		router.ServeHTTP(w, r)

		if mock.key != "key2" {
			t.Errorf("Expected key2 but got %s", mock.key)
		}
	})

	t.Run("UnHappy Days", func(t *testing.T) {
		rf, mock := initRESTLayers(t, false, false)
		w := httptest.NewRecorder()
		data := "{\"key2\": { \"name\": \"name2\",\"city\": \"city2\"}"
		r := httptest.NewRequest("POST", "/port", strings.NewReader(data))
		r.Header.Add("Content-Type", "application/json")
		router := mux.NewRouter()
		router.HandleFunc("/port", rf.PostHandler).Methods("POST")
		router.ServeHTTP(w, r)

		if mock.key != "" {
			t.Errorf("Expected nothing but got %s", mock.key)
		}
	})
}

func TestFramework_RespondToHeartBeat(t *testing.T) {
	rf, _ := initRESTLayers(t, false, true)
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		rf   Framework
		args args
	}{
		{
			"Happy Days",
			rf,
			args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/heartbeat", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rf.RespondToHeartBeat(tt.args.w, tt.args.r)
		})
	}
}
