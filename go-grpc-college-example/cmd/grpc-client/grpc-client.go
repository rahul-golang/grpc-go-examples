package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	v1 "github.com/rahul/go-grpc-college-example/pkg/api/v1"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

func main() {
	// get configuration
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	StudentClient := v1.NewStudentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call Create
	req := v1.CreateStudentReq{
		Api: &v1.API{Api: "v1"},
		Student: &v1.Student{
			RollNo:     1,
			Name:       "rahul",
			Department: "scirnce",
		},
	}
	res, err := StudentClient.CreateNewStudent(ctx, &req)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", res)

	req1 := &v1.ReadAllStudentRecordsReq{
		Api: &v1.API{Api: "v1"},
	}
	stream, err := StudentClient.ReadAllStudentRecords(ctx)

	waitc := make(chan interface{})

	//msg := &v1.ReadAllStudentRecordsReq{"sup"}

	go func() {
		for {
			log.Println("Sleeping...")
			time.Sleep(2 * time.Second)
			log.Println("Sending msg...")
			stream.Send(req1)

		}
	}()

	go func() {
		for {
			log.Println("Sleeping...")
			time.Sleep(2 * time.Second)
			log.Println("Recive msg...")
			aa, errr := stream.Recv()
			if errr != nil {
				log.Println(errr)
			}

			log.Println(aa)
			waitc <- aa
		}
	}()
	<-waitc
	stream.CloseSend()

	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", stream)

}
