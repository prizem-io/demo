package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/prizem-io/demo/pkg/propagate"
)

const (
	envBackendServiceURI = "BACKEND_URI"
)

type messageClient struct{}

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		ctx := propagate.Incoming(r.Context(), r.Header)

		username, password, ok := r.BasicAuth()
		if !ok || !(username == "demo" && password == "demo") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.Printf("[INFO] %s %s %s", r.Method, r.URL.Path, r.UserAgent())

		host := os.Getenv(envBackendServiceURI)
		if len(host) == 0 {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("No backend service is working\n"))
			return
		}

		req, err := http.NewRequest("GET", host, nil)
		if err != nil {
			log.Printf("[ERROR] Failed to create new request: %s", err)
			return
		}
		propagate.Outgoing(ctx, req.Header)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("[ERROR] Failed to do request: %s", err)
			return
		}

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	log.Printf("[INFO] Started server on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("[ERROR] Failed to start server: %s", err)
	}
}
