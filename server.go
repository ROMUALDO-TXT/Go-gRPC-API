package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "github.com/ROMUALDO-TXT/Go-gRPC-API/management"
	"google.golang.org/grpc"
)

const (
	port = ":5005"
)

type ManagementServer struct {
	pb.UnimplementedManagementServer
}

func (s *ManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received %v", in.GetName())
	var user_id int32 = int32(rand.Intn(1000))

	return &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("fail to listen to port: %v", port)
	}

	s := grpc.NewServer()

	pb.RegisterManagementServer(s, &ManagementServer{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}
