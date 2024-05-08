package main

import (
	"context"
	"io"
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

func (s *server) StreamPingPong(stream pb.PingPongService_StreamPingPongServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// ストリームの終了
			return nil
		}
		if err != nil {
			log.Fatalf("Failed to receive a note : %v", err)
			return err
		}
		// 受け取った Ping メッセージに応じて Pong を返す
		if err := stream.Send(&pb.Pong{Message: "Pong received: " + in.Message}); err != nil {
			log.Fatalf("Failed to send a pong : %v", err)
			return err
		}
	}
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
