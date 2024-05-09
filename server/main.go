package main

import (
	"context"
	"fmt"
	pb "github.com/pikachu0310/grpc_test/server/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
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

func (s *server) ReceivePongStream(req *pb.Empty, stream pb.PingPongService_ReceivePongStreamServer) error {
	for {
		// ここでは例として、1秒ごとにPongメッセージを送信
		time.Sleep(1 * time.Second)
		if err := stream.Send(&pb.Pong{Message: "Pong from stream"}); err != nil {
			log.Fatalf("Failed to send pong stream: %v", err)
			return err
		}
	}
}

func (s *server) PingAndStreamPong(req *pb.Ping, stream pb.PingPongService_PingAndStreamPongServer) error {
	for i := 0; i < 10; i++ { // ここでは例として、10回Pongメッセージを送信
		time.Sleep(500 * time.Millisecond) // 少し間隔を開ける
		pongMsg := fmt.Sprintf("Pong for %s", req.Message)
		if err := stream.Send(&pb.Pong{Message: pongMsg}); err != nil {
			log.Fatalf("Failed to send pong for ping: %v", err)
			return err
		}
	}
	return nil
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
