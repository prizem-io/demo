package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc"

	"github.com/prizem-io/demo/pkg/propagate"
	pb "github.com/prizem-io/demo/proto/message"
)

const (
	envMessageServiceURI     = "MESSAGE_URI"
	envGRPCMessageServiceURI = "GRPC_MESSAGE_URI"
)

func main() {
	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		ctx := propagate.Incoming(r.Context(), r.Header)
		log.Printf("[INFO] %s %s %s", r.Method, r.URL.Path, r.UserAgent())

		host := os.Getenv(envMessageServiceURI)
		grpcHost := os.Getenv(envGRPCMessageServiceURI)

		if len(host) != 0 {
			// Do normal http request
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

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
		} else if len(grpcHost) != 0 {
			// Do gRPC request
			conn, err := grpc.Dial(grpcHost, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("[ERROR] Faield to connect: %v", err)
				return
			}
			defer conn.Close()

			client := pb.NewMessageClient(conn)

			ctx := context.Background()
			resp, err := client.Hello(ctx, &pb.HelloRequest{
				Name: "demo",
			})
			if err != nil {
				log.Fatalf("[ERROR] Faield: %v", err)
				return
			}

			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(resp.Message + "\n"))

			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("No backend message service is working\n"))

			return
		}
	})

	log.Printf("[INFO] Started server on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("[ERROR] Failed to start server: %s", err)
	}
}
