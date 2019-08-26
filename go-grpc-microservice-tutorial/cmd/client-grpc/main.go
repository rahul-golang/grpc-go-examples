package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	v1 "github.com/rahul/go-grpc-microservice-tutorial/pkg/api/v1"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

type loginCreds struct {
	Username, Password string
}

func (c *loginCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"username": c.Username,
		"password": c.Password,
	}, nil
}

func (c *loginCreds) RequireTransportSecurity() bool {
	return true
}

func main() {
	// get configuration
	address := flag.String("server", "localhost:9090", "gRPC server in format host:port")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

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
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", res1)

	BookClient := v1.NewBookServiceClient(conn)

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	// Call Create
	BookCreateReq := v1.Book{
		Api:        apiVersion,
		BookID:     0,
		BookName:   "GRPC Book",
		BookAuthor: "Rahul",
		Pages:      100,
	}
	BookCreateResp, BookCreateErr := BookClient.CreateBookRecord(ctx, &BookCreateReq)
	if BookCreateErr != nil {
		log.Fatalf("Create failed: %v", BookCreateErr)
	}
	log.Printf("Create result: <%+v>\n\n", BookCreateResp)

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
