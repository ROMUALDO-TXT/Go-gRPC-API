package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"

	pb "github.com/ROMUALDO-TXT/Go-gRPC-API/management"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	port = ":5005"
)

func NewManagementServer() *ManagementServer {
	return &ManagementServer{}
}

type ManagementServer struct {
	pb.UnimplementedManagementServer
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
	readBytes, err := ioutil.ReadFile("./temp/users.json")
	var users_list *pb.UserList = &pb.UserList{}
	var user_id int32 = int32(rand.Intn(1000))
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}

	if err != nil {
		if os.IsNotExist(err) {
			log.Print("file not found. Creating a new file")
			users_list.Users = append(users_list.Users, created_user)
			jsonBytes, err := protojson.Marshal(users_list)
			if err != nil {
				log.Fatalf("JSON Marshaling failed: %v", err)
			}
			if err := ioutil.WriteFile("./temp/users.json", jsonBytes, 0664); err != nil {
				log.Fatalf("failed to write file: %v", err)
			}
			return created_user, nil
		} else {
			log.Fatalf("error reading file: %v", err)
		}
	}

	if err := protojson.Unmarshal(readBytes, users_list); err != nil {
		log.Fatalf("failed to parse user list: %v", err)
	}
	users_list.Users = append(users_list.Users, created_user)

	jsonBytes, err := protojson.Marshal(users_list)
	if err != nil {
		log.Fatalf("JSON Marshaling failed: %v", err)
	}
	if err := ioutil.WriteFile("./temp/users.json", jsonBytes, 0664); err != nil {
		log.Fatalf("failed to write file: %v", err)
	}

	return created_user, nil
}

func (s *ManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	jsonBytes, err := ioutil.ReadFile("./temp/users.json")
	if err != nil {
		log.Fatalf("failed to read from file: %v", err)
	}
	var users_list *pb.UserList = &pb.UserList{}

	if err := protojson.Unmarshal(jsonBytes, users_list); err != nil {
		log.Fatalf("Unmarshaling failed: %v", err)
	}

	return users_list, nil
}

func main() {
	var mgmt_server *ManagementServer = NewManagementServer()

	if err := mgmt_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
