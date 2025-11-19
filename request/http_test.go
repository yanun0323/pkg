package request

import (
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"testing"
)

var mockServerOn atomic.Bool

func StartMockServer() {
	if mockServerOn.Swap(true) {
		return
	}

	mux := http.DefaultServeMux
	mux.HandleFunc("GET /good", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("{\"status\":\"good\"}"))
	})

	mux.HandleFunc("POST /good", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			_, _ = fmt.Fprintf(w, "{\"status\":\"bad\",\"error\":\"%s\"}", err)
			return
		}

		_, _ = fmt.Fprintf(w, "{\"status\":\"good\",\"body\":\"%s\"}", string(body))
	})

	mux.HandleFunc("GET /bad", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"status\":\"bad\"}"))
	})

	mux.HandleFunc("POST /bad", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"status\":\"bad\"}"))
	})

	go http.ListenAndServe("0.0.0.0:8888", mux)
}

// go test -bench=. -benchmem ./...
func BenchmarkRequest(b *testing.B) {
	StartMockServer()
	b.Run("Get - Good", func(b *testing.B) {
		for b.Loop() {
			res, err := New(http.MethodGet, "http://127.0.0.1:8888/good").
				WithContext(b.Context()).
				WithHeader("HEADER", "VALUE").
				WithQueryParam("query", "value").
				WithQueryParam("query2", "value2").
				WithHook(func(r Request) error {
					return nil
				}).Send(http.DefaultClient.Do)
			if err != nil {
				b.Fatalf("send request, err: %s\n", err)
			}

			var response map[string]string
			if err := res.WithCheckStatus().Decode(&response); err != nil {
				b.Fatalf("decode response, err: %s\n", err)
			}

			if len(response["status"]) == 0 || response["status"] != "good" {
				b.Fatalf("response should NOT be empty, status=%s", response["status"])
			}
		}
	})
}
