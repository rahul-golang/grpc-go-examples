package v1

// import (
// 	"context"
// 	"database/sql"

// 	v1 "github.com/rahul/go-grpc-microservice-tutorial/pkg/api/v1"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// const (
// 	// apiVersion is version of API is provided by server
// 	apiVersion = "v1"
// )

// // DatabaseConnection is implementation of v1.ToDoServiceServer proto interface
// type DatabaseConnection struct {
// 	db *sql.DB
// }

// // CreatedDatabaseConnection creates ToDo service
// func CreatedDatabaseConnection(db *sql.DB) v1.ToDoServiceServer {
// 	return &DatabaseConnection{db: db}
// }

// // checkAPI checks if the API version requested by client is supported by server
// func (s *DatabaseConnection) checkAPI(api string) error {
// 	// API version is "" means use current version of the service
// 	if len(api) > 0 {
// 		if apiVersion != api {
// 			return status.Errorf(codes.Unimplemented,
// 				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
// 		}
// 	}
// 	return nil
// }

// // connect returns SQL database connection from the pool
// func (s *DatabaseConnection) connect(ctx context.Context) (*sql.Conn, error) {
// 	c, err := s.db.Conn(ctx)
// 	if err != nil {
// 		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
// 	}
// 	return c, nil
// }

// // func (s *DatabaseConnection) Create(context.Context, *v1.CreateRequest) (*v1.CreateResponse, error) {
// // 	return &v1.CreateResponse{}, nil

// //}
