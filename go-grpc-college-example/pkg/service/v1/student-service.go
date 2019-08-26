package v1

import (
	"context"
	"database/sql"
	"io"
	"log"

	v1 "github.com/rahul/go-grpc-college-example/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	//APIVersion  version
	APIVersion = "v1"
)

//StudentSevicesStruct is Initialize with service stuct
type StudentSevicesStruct struct {
	DB *sql.DB
}

// checkAPI checks if the API version requested by client is supported by server
func (studentSevicesStruct *StudentSevicesStruct) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if APIVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", APIVersion, api)
		}
	}
	return nil
}

// connect returns SQL database connection from the pool
func (studentSevicesStruct *StudentSevicesStruct) connect(ctx context.Context) (*sql.Conn, error) {
	/*studentSevicesStruct.DB.Conn(ctx) represents a single database connection rather than a pool of database connections. Prefer running queries from DB unless there is a specific need for a continuous single database connection.
	A Conn must call Close to return the connection to the database pool and may do so concurrently with a running query.
	After a call to Close, all operations on the connection fail with ErrConnDone.
	*/
	c, err := studentSevicesStruct.DB.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// NewStudentSevicesStruct initialize stuct
func NewStudentSevicesStruct(db *sql.DB) v1.StudentServiceServer {

	return &StudentSevicesStruct{DB: db}
}

// IntializeAPIVersion initialize api version
func IntializeAPIVersion() *v1.API {
	return &v1.API{Api: APIVersion}
}

// IntializeStudent student with proto Struct
func IntializeStudent(RollNo int64, Name string) *v1.Student {
	return &v1.Student{RollNo: RollNo,
		Name: Name}
}

// IntializeStudentArray with multiple students
func IntializeStudentArray() *v1.ReadAllStudentRecordsResp {
	var stud []*v1.Student
	stud = append(stud, IntializeStudent(1, "rahul"))
	stud = append(stud, IntializeStudent(2, "rahul"))
	return &v1.ReadAllStudentRecordsResp{Api: IntializeAPIVersion(), Student: stud}
}

//CreateNewStudent students in database
func (studentSevicesStruct *StudentSevicesStruct) CreateNewStudent(ctx context.Context, Request *v1.CreateStudentReq) (*v1.CreateStudentResp, error) {

	if err := studentSevicesStruct.checkAPI(Request.Api.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := studentSevicesStruct.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// insert into Student data entity data
	QueryResponse, Err := c.ExecContext(ctx, "INSERT INTO student(`name`,`department`) VALUES(?, ?)",
		Request.Student.Name, Request.Student.Department)
	if Err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into student-> "+Err.Error())
	}

	// get ID of creates student
	ID, Err := QueryResponse.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created ToDo-> "+err.Error())
	}

	return &v1.CreateStudentResp{Api: IntializeAPIVersion(),
		Student: IntializeStudent(ID, Request.Student.Name),
	}, nil
}

// func (studentSevicesStruct *StudentSevicesStruct) ReadAllStudentRecords(ctx context.Context, req *v1.ReadAllStudentRecordsReq) (*v1.ReadAllStudentRecordsResp, error) {

// 	return &v1.ReadAllStudentRecordsResp{Api: IntializeApiVersion(),
// 		Student: IntializeStudentArray()}, nil
// }

// ReadAllStudentRecords Streming
func (studentSevicesStruct *StudentSevicesStruct) ReadAllStudentRecords(stream v1.StudentService_ReadAllStudentRecordsServer) error {

	for {
		in, err := stream.Recv()
		log.Println("Received value")
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		stream.Send(IntializeStudentArray())
		log.Println(in)

	}

	// return &v1.ReadAllStudentRecordsResp{Api: IntializeApiVersion(),
	// 	Student: IntializeStudentArray()}, nil
}
