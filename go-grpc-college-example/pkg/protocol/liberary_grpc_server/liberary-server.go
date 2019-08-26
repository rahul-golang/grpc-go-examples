package liberary_grpc_server

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"

	v1 "github.com/rahul/go-grpc-college-example/pkg/api/v1"
	service_v1 "github.com/rahul/go-grpc-college-example/pkg/service/v1"
	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, db *sql.DB, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	v1.RegisterStudentServiceServer(server, service_v1.NewStudentSevicesStruct(db))

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}
