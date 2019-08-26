package v1

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	v1 "github.com/rahul/go-grpc-microservice-http-tutorial/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// bookServiceServer is implementation of v1.BookServiceServer proto interface
type bookServiceServer struct {
	db *sql.DB
}

// NewBookServiceServer creates Book service
func NewBookServiceServer(db *sql.DB) v1.BookServiceServer {
	return &bookServiceServer{db: db}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *bookServiceServer) checkAPI(api string) error {
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
func (s *bookServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// Create new todo task
func (s *bookServiceServer) CreateBookRecord(ctx context.Context, req *v1.Book) (*v1.Book, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// insert Book entity data
	res, err := c.ExecContext(ctx, "INSERT INTO Book(`BookName`, `BookAuthor`, `Pages`) VALUES(?, ?, ?)",
		req.BookID, req.BookName, req.Pages)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into BOOK-> "+err.Error())
	}

	// get ID of creates Book
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created ToDo-> "+err.Error())
	}

	//88888888888888888888888
	//	address1 := flag.String("server", "localhost:9090", "gRPC server in format host:port")
	//	flag.Parse()

	// Set up a connection to the server.
	conn2, err2 := grpc.Dial("localhost:9090", grpc.WithInsecure()) // "(*address1"
	if err2 != nil {
		log.Fatalf("did not connect: %v", err2)
	}
	defer conn2.Close()

	New_library_cli := v1.NewLibraryServiceClient(conn2)

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	// Call Create
	req11 := v1.IssueBookReq{
		Api: apiVersion,
		Book: &v1.Book{
			Api:        apiVersion,
			BookID:     2,
			BookName:   "222",
			BookAuthor: "AAA",
			Pages:      1,
		},
		IssuerName: "rahul",
	}
	fmt.Println("connection open", req11)
	res113, err3 := New_library_cli.IssueBook(ctx, &req11) //Create(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err3)
	}
	log.Printf("Create result: <%+v>\n\n", res113)

	// Set up a connection to the server.

	return &v1.Book{
		Api:        apiVersion,
		BookID:     id,
		BookName:   req.BookName,
		BookAuthor: req.BookAuthor,
		Pages:      req.Pages,
	}, nil
}

// Read todo task
func (s *bookServiceServer) ReadBookRecord(ctx context.Context, req *v1.ReadBookRecordReq) (*v1.Book, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// query ToDo by ID
	rows, err := c.QueryContext(ctx, "SELECT `BookId`, `BookName`, `BookAuthor`, `Pages` FROM Book WHERE `BookId`=?",
		req.BookID)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo-> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from ToDo-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found",
			req.BookID))
	}

	// get ToDo data
	var td v1.Book
	//var reminder time.Time
	if err := rows.Scan(&td.BookID, &td.BookName, &td.BookAuthor, &td.Pages); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from ToDo row-> "+err.Error())
	}
	// td.Reminder, err = ptypes.TimestampProto(reminder)
	// if err != nil {
	// 	return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
	// }

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple ToDo rows with ID='%d'",
			req.BookID))
	}

	return &td, nil

}

// Update todo task
func (s *bookServiceServer) UpdateBookRecord(ctx context.Context, req *v1.Book) (*v1.Book, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// reminder, err := ptypes.Timestamp(req.ToDo.Reminder)
	// if err != nil {
	// 	return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
	// }

	// update ToDo
	res, err := c.ExecContext(ctx, "UPDATE Book SET `BookName`=?, `BookAuthor`=?, `Pages`=? WHERE `ID`=?",
		req.BookName, req.BookAuthor, req.Pages, req.BookID)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update ToDo-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found",
			req.BookID))
	}

	return &v1.Book{
		Api:        apiVersion,
		BookID:     req.BookID,
		BookName:   req.BookName,
		BookAuthor: req.BookAuthor,
		Pages:      req.Pages,
	}, nil
}

// Delete todo task
func (s *bookServiceServer) DeleteBookRecord(ctx context.Context, req *v1.Book) (*v1.Book, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// delete ToDo
	res, err := c.ExecContext(ctx, "DELETE FROM ToDo WHERE `ID`=?", req.BookID)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete ToDo-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found",
			req.BookID))
	}

	return &v1.Book{
		Api:        apiVersion,
		BookID:     req.BookID,
		BookName:   req.BookName,
		BookAuthor: req.BookAuthor,
		Pages:      req.Pages,
	}, nil
}

// Read all todo tasks
func (s *bookServiceServer) ReadAllBook(ctx context.Context, req *v1.ReadAllBookReq) (*v1.ReadAllBookResp, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// get ToDo list
	rows, err := c.QueryContext(ctx, "SELECT `BookId`, `BookName`, `BookAuthor`, `Pages` FROM Book")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo-> "+err.Error())
	}
	defer rows.Close()

	//var reminder time.Time
	list := []*v1.Book{}
	for rows.Next() {
		td := new(v1.Book)
		if err := rows.Scan(&td.BookID, &td.BookName, &td.BookAuthor, &td.Pages); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from ToDo row-> "+err.Error())
		}
		// td.Reminder, err = ptypes.TimestampProto(reminder)
		// if err != nil {
		// 	return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
		// }
		list = append(list, td)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from ToDo-> "+err.Error())
	}

	return &v1.ReadAllBookResp{
		Api:  apiVersion,
		Book: list,
	}, nil
}
