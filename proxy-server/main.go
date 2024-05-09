package main

import (
	"context"
	pingpong "github.com/pikachu0310/grpc_test/server/proto"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	http.HandleFunc("/", handler)
	log.Println("Starting proxy server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	// Determine which service to route based on URL
	serviceMethod := strings.TrimPrefix(r.URL.Path, "/")
	switch serviceMethod {
	case "SendPing":
		handleUnarySendPing(w, r)
	case "ReceivePongStream":
		handleServerStreamingReceivePongStream(w, r)
	case "PingAndStreamPong":
		handleServerStreamingPingAndStreamPong(w, r)
	default:
		http.Error(w, "Method not supported", http.StatusNotFound)
	}
}

func handleUnarySendPing(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "Failed to connect to gRPC server: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pingpong.NewPingPongServiceClient(conn)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req pingpong.Ping
	if err := proto.Unmarshal(body, &req); err != nil {
		http.Error(w, "Failed to parse request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := client.SendPing(context.Background(), &req)
	if err != nil {
		http.Error(w, "gRPC server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := proto.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/grpc-web+proto")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleServerStreamingReceivePongStream(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "Failed to connect to gRPC server: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pingpong.NewPingPongServiceClient(conn)
	stream, err := client.ReceivePongStream(context.Background(), &pingpong.Empty{})
	if err != nil {
		http.Error(w, "gRPC server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stream.CloseSend()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/grpc-web+proto")
	w.WriteHeader(http.StatusOK)

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Failed to receive from stream: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := proto.Marshal(resp)
		if err != nil {
			http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(data)
		flusher.Flush() // Flush the data immediately
	}
}

func handleServerStreamingPingAndStreamPong(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "Failed to connect to gRPC server: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pingpong.NewPingPongServiceClient(conn)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req pingpong.Ping
	if err := proto.Unmarshal(body, &req); err != nil {
		http.Error(w, "Failed to parse request: "+err.Error(), http.StatusBadRequest)
		return
	}

	stream, err := client.PingAndStreamPong(context.Background(), &req)
	if err != nil {
		http.Error(w, "gRPC server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stream.CloseSend()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/grpc-web+proto")
	w.WriteHeader(http.StatusOK)

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Failed to receive from stream: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := proto.Marshal(resp)
		if err != nil {
			http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(data)
		flusher.Flush() // Flush the data immediately
	}
}
