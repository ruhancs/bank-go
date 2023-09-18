package grpcapi

import (
	"context"

	"github.com/lib/pq"
	db "github.com/ruhancs/bank-go/db/sqlc"
	"github.com/ruhancs/bank-go/pb"
	"github.com/ruhancs/bank-go/util"
	"github.com/ruhancs/bank-go/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil,invalidArgumentError(violations)
	}

	hashedPassword,err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	arg := db.CreateUserParams{
		Username: req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName: req.GetFullname(),
		Email: req.GetEmail(),
	}

	user,err := server.store.CreateUser(ctx,arg)
	if err != nil {
		//erro ao criar user, verificar o tipo do db
		if pqErr,ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation": //verificar erro de username ja existe
				return nil, status.Errorf(codes.AlreadyExists, "username alredy registered")
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	resp := &pb.CreateUserResponse{
		User: converterUser(user),
	}

	return resp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := val.ValidateFullName(req.GetFullname()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}