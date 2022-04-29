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

func NewManagementServer() *ManagementServer {
	return &ManagementServer{
		user_list: &pb.UserList{},
	}
}

type ManagementServer struct {
	pb.UnimplementedManagementServer
	user_list *pb.UserList
}

func (server *ManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("fail to listen to port: %v", port)

	}

	s := grpc.NewServer()

	pb.RegisterManagementServer(s, server)

	log.Printf("server listening at %v", lis.Addr())

	return s.Serve(lis)
}

func (s *ManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received %v", in.GetName())
	var user_id int32 = int32(rand.Intn(1000))
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}
	s.user_list.Users = append(s.user_list.Users, created_user)
	return created_user, nil
}

func (s *ManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	return s.user_list, nil
}

func main() {
	var mgmt_server *ManagementServer = NewManagementServer()

	if err := mgmt_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
