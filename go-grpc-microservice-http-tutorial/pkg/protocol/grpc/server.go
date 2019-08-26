package grpc

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"

	v1 "github.com/rahul/go-grpc-microservice-http-tutorial/pkg/api/v1"
	ser_v1 "github.com/rahul/go-grpc-microservice-http-tutorial/pkg/service/v1"
	"google.golang.org/grpc"
)

// RunServer runs gRPC service to publish ToDo service // v1API1 v1.BookServiceServer
/*func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, v1APIBook v1.BookServiceServer, port string) error {
 */
func RunServer(ctx context.Context, db *sql.DB, port string) error {

	// // register service
	server1 := grpc.NewServer()
	//server2 := grpc.NewServer()

	v1.RegisterToDoServiceServer(server1, ser_v1.NewToDoServiceServer(db))
	v1.RegisterBookServiceServer(server1, ser_v1.NewBookServiceServer(db))
	v1.RegisterLibraryServiceServer(server1, ser_v1.NewLibraryServiceServer(db))

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server1.GracefulStop()
			//			server2.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")

	listen1, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// port2 := "8000"
	// listen2, err2 := net.Listen("tcp", "localhost:"+port2)
	// if err2 != nil {
	// 	return err2
	// }

	ReturnError := server1.Serve(listen1)
	if ReturnError != nil {
		return ReturnError
	}

	// ReturnError = server2.Serve(listen2)
	// if ReturnError != nil {
	// 	return ReturnError
	// }

	return ReturnError
}
