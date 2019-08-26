package v1

import (
	"context"
	"database/sql"
	"fmt"

	v1 "github.com/rahul/grpc-assignments/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

// userManagementService is implementation of v1.UserManagmentServer proto interface
type userManagementService struct {
	db *sql.DB
}

// NewUserManagementService creates UserManagementService service
func NewUserManagementService(db *sql.DB) v1.UserManagmentServer {
	return &userManagementService{db: db}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *userManagementService) checkAPI(api string) error {
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
func (s *userManagementService) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// Create new User
func (s *userManagementService) CreateUser(ctx context.Context, req *v1.User) (*v1.CreateUserResp, error) {
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

	// insert user entity data
	res, err := c.ExecContext(ctx, "INSERT INTO user(`name`, `email`, `username`,`qualification`,`experience`,`password`,`invitation_flag`) VALUES(?, ?, ?, ? ,? ,? ,?)",
		req.Name, req.Email, req.UserName, req.Qualification, req.Experience, req.Password, req.InvitationFlag)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into user-> "+err.Error())
	}

	// get UserId of creates User
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created user-> "+err.Error())
	}

	return &v1.CreateUserResp{
		Api:      apiVersion,
		UserId:   id,
		RespCode: "ssuccess",
		Message:  "user created sucessfully",
	}, nil
}

// Read todo task
func (s *userManagementService) GetUser(ctx context.Context, req *v1.GetUserReq) (*v1.User, error) {
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
	rows, err := c.QueryContext(ctx, "SELECT `name`, `email`, `username`,`qualification`,`experience`,`password`,`invitation_flag` FROM user WHERE `username`=? and `password`=?",
		req.UserName, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to Get user from user-> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from user-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("User with username and Password is not found"))
	}

	// get User data
	var td v1.User

	if err := rows.Scan(&td.Name, &td.Email, &td.UserName, &td.Qualification, &td.Experience, &td.Password, &td.InvitationFlag); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from user row-> "+err.Error())
	}

	// if err := rows.Scan(&td.Name, &td.Email, &td.UserName, &td.Qualification,&td.Experience,&td.Password,&td.InvitationFlag); err != nil {
	// 	return nil, status.Error(codes.Unknown, "failed to retrieve field values from user row-> "+err.Error())
	// }

	// if rows.Next() {
	// 	return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple user rows with username"))
	// }

	return &td, nil

}
