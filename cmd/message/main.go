package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/sayHello", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO] %s %s %s", r.Method, r.URL.Path, r.UserAgent())

		response := struct {
			Message string `json:"message"`
		}{
			Message: "hello, world",
		}

		w.Header().Add("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(&response); err != nil {
			log.Printf("[ERROR] Failed to encode message: %s", err)
			return
		}
	})

	log.Printf("[INFO] Started server on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("[ERROR] Failed to start server")
	}
}
