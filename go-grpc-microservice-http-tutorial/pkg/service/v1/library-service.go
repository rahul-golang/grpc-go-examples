package v1

import (
	"context"
	"database/sql"
	"fmt"

	v1 "github.com/rahul/go-grpc-microservice-http-tutorial/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// toDoServiceServer is implementation of v1.ToDoServiceServer proto interface
type libraryServiceServer struct {
	db *sql.DB
}

// NewLibraryServiceServer creates ToDo service
func NewLibraryServiceServer(db *sql.DB) v1.LibraryServiceServer {
	return &libraryServiceServer{db: db}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *libraryServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

// connect returns SQL database connection from the pool
func (s *libraryServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// Create new todo task
func (s *libraryServiceServer) IssueBook(ctx context.Context, req *v1.IssueBookReq) (*v1.IssueBookResp, error) {
	// check if the API version requested by client is supported by server
	fmt.Println("called")
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// insert ToDo entity data
	res, err := c.ExecContext(ctx, "INSERT INTO Library(`BookID`, `BookName`,`BookAuthor`,`Pages`,`IssuerName`) VALUES(?, ?, ?,?,?)",
		req.Book.BookID, req.Book.BookName, req.Book.BookAuthor, req.Book.Pages, req.IssuerName)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into Library-> "+err.Error())
		fmt.Println(err.Error())
	}

	fmt.Println("called")
	// get ID of creates ToDo
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created Library-> "+err.Error())
	}
	req.Book.BookID = id

	return &v1.IssueBookResp{
		Api:        apiVersion,
		Book:       req.Book,
		IssuerName: req.IssuerName,
		Status:     "ISUUED",
	}, nil
}
