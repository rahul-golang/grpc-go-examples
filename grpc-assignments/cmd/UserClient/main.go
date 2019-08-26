package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	v1 "github.com/rahul/grpc-assignments/pkg/api/v1"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

func main() {
	// get configuration
	address := flag.String("server", "localhost:9090", "gRPC server in format host:port")
	flag.Parse()

	// Set up a connection to the server.
	// Dial creates a client connection to the given target

	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//initialize client connection struct in protobuf
	c := v1.NewUserManagmentClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call Create
	req1 := v1.User{
		Api:            apiVersion,
		Name:           "Rahul",
		Email:          "r@example.com",
		UserName:       "rrrr",
		Qualification:  "AAA",
		Experience:     1.0,
		Password:       "rrrr",
		InvitationFlag: 0,
	}
	res1, err := c.CreateUser(ctx, &req1)
	if err != nil {
		log.Fatalf("Create  User failed: %v", err)
	}
	log.Printf("Create User result: <%+v>\n\n", res1)

	getUser := v1.GetUserReq{
		Api:      apiVersion,
		UserName: "rrrr",
		Password: "rrrr",
	}

	resp2, err := c.GetUser(ctx, &getUser)
	if err != nil {
		log.Fatalf("GetUser Error failed :%v ", err)
	}
	log.Printf("Get User result: <%+v>\n\n", resp2)

}
