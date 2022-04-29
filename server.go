package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/ROMUALDO-TXT/Go-gRPC-API/management"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

func NewManagementServer() *ManagementServer {
	return &ManagementServer{}
}

type ManagementServer struct {
	conn *pgx.Conn
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

func (server *ManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received %v", in.GetName())

	createSql := `
		create table if it not exists users(
			id SERIAL PRIMARY KEY
			name varchar(200)
			age int
		);
		`

	_, err := server.conn.Exec(context.Background(), createSql)

	if err != nil {
		fmt.Fprint(os.Stderr, "Table Creation Failed: %v\n", err)
		os.Exit(1)
	}

	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge()}
	tx, err := server.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed: %v", err)
	}
	_, err = tx.Exec(context.Background(), "insert into users(name, age) values ($1,$2)", created_user.Name, created_user.Age)

	if err != nil {
		log.Fatalf("tx.Exec failed: %v", err)
	}

	tx.Commit(context.Background())

	return created_user, nil
}

func (server *ManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {

	var users_list *pb.UserList = &pb.UserList{}
	rows, err := server.conn.Query(context.Background(), "select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := pb.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		users_list.Users = append(users_list.Users, &user)
	}
	return users_list, nil
}

func main() {
	var db_url string = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	conn, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer conn.Close(context.Background())
	var mgmt_server *ManagementServer = NewManagementServer()
	mgmt_server.conn = conn

	if err := mgmt_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
