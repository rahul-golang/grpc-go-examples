package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	v1 "github.com/rahul/go-grpc-microservice-http-tutorial/pkg/api/v1"
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
	c := v1.NewToDoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(t)
	pfx := t.Format(time.RFC3339Nano)

	// Call Create
	req1 := v1.CreateRequest{
		Api: apiVersion,
		ToDo: &v1.ToDo{
			Title:       "title (" + pfx + ")",
			Description: "description (" + pfx + ")",
			Reminder:    reminder,
		},
	}
	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatalf("Create  TODO failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", res1)

	BookClient := v1.NewBookServiceClient(conn)
	BookClientReq := v1.Book{
		Api:        apiVersion,
		BookName:   "The Earth",
		BookAuthor: "shiva",
		Pages:      123,
	}

	BookClientResp, err := BookClient.CreateBookRecord(ctx, &BookClientReq)

	if err != nil {
		log.Fatalf("Creating Book Records : %v", err)
	}
	log.Printf("Create result: %+v \n\n", BookClientResp)

	//BookClient := v1.NewBookServiceClient(conn)
	ReadAllBookReq := v1.ReadAllBookReq{
		Api: apiVersion,
	}

	ReadAllBookResp, err := BookClient.ReadAllBook(ctx, &ReadAllBookReq)

	if err != nil {
		log.Fatalf("Creating Book Records : %v", err)
	}
	log.Printf("Create result:<%+v>\n\n", ReadAllBookResp)

	// address1 := flag.String("server1", "localhost:6060", "gRPC server in format host:port")
	// flag.Parse()

	// // Set up a connection to the server.
	// conn2, err2 := grpc.Dial(*address1, grpc.WithInsecure())
	// if err2 != nil {
	// 	log.Fatalf("did not connect: %v", err2)
	// }
	// defer conn2.Close()

	// New_library_cli := v1.NewLibraryServiceClient(conn2)

	// //ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// //defer cancel()

	// // Call Create
	// req11 := v1.IssueBookReq{
	// 	Api: apiVersion,
	// 	Book: &v1.Book{
	// 		Api:        apiVersion,
	// 		BookID:     2,
	// 		BookName:   "222",
	// 		BookAuthor: "AAA",
	// 		Pages:      1,
	// 	},
	// 	IssuerName: "rahul",
	// }
	// fmt.Println("connection open", req11)
	// res113, err3 := New_library_cli.IssueBook(ctx, &req11) //Create(ctx, &req1)
	// if err != nil {
	// 	log.Fatalf("Create failed: %v", err3)
	// }
	// log.Printf("Create result: <%+v>\n\n", res113)

	// id := res1.Id

	// // Read
	// req2 := v1.ReadRequest{
	// 	Api: apiVersion,
	// 	Id:  id,
	// }
	// res2, err := c.Read(ctx, &req2)
	// if err != nil {
	// 	log.Fatalf("Read failed: %v", err)
	// }
	// log.Printf("Read result: <%+v>\n\n", res2)

	// // Update
	// req3 := v1.UpdateRequest{
	// 	Api: apiVersion,
	// 	ToDo: &v1.ToDo{
	// 		Id:          res2.ToDo.Id,
	// 		Title:       res2.ToDo.Title,
	// 		Description: res2.ToDo.Description + " + updated",
	// 		Reminder:    res2.ToDo.Reminder,
	// 	},
	// }
	// res3, err := c.Update(ctx, &req3)
	// if err != nil {
	// 	log.Fatalf("Update failed: %v", err)
	// }
	// log.Printf("Update result: <%+v>\n\n", res3)

	// // Call ReadAll
	// req4 := v1.ReadAllRequest{
	// 	Api: apiVersion,
	// }
	// res4, err := c.ReadAll(ctx, &req4)
	// if err != nil {
	// 	log.Fatalf("ReadAll failed: %v", err)
	// }
	// log.Printf("ReadAll result: <%+v>\n\n", res4)

	// // Delete
	// req5 := v1.DeleteRequest{
	// 	Api: apiVersion,
	// 	Id:  id,
	// }
	// res5, err := c.Delete(ctx, &req5)
	// if err != nil {
	// 	log.Fatalf("Delete failed: %v", err)
	// }
	// log.Printf("Delete result: <%+v>\n\n", res5)

}
