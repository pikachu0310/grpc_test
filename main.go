package main

import (
	"context"
	"log"
	"net"

	pb "github.com/pikachu0310/grpc_test/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPingPongServiceServer
}

func (s *server) SendPing(ctx context.Context, in *pb.Ping) (*pb.Pong, error) {
	return &pb.Pong{Message: "pong"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPingPongServiceServer(s, &server{})
	log.Println("Server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
