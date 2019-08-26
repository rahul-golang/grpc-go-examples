package grpc

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"

	v1 "github.com/rahul/grpc-assignments/pkg/api/v1"
	ser_v1 "github.com/rahul/grpc-assignments/pkg/service/v1"
	"google.golang.org/grpc"
)

// RunServer runs gRPC service to publish ToDo service // v1API1 v1.BookServiceServer
/*func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, v1APIBook v1.BookServiceServer, port string) error {
 */
func RunServer(ctx context.Context, db *sql.DB, port string) error {

	// register service
	server := grpc.NewServer()

	v1.RegisterUserManagmentServer(server, ser_v1.NewUserManagementService(db))

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()
			//			server2.GracefulStop()

			<-ctx.Done()
		}
	}()
	// start gRPC server
	log.Println("starting gRPC server...")

	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	ReturnError := server.Serve(listen)
	if ReturnError != nil {
		return ReturnError
	}
	return ReturnError
}
