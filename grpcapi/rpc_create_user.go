package grpcapi

import (
	"context"

	"github.com/lib/pq"
	db "github.com/ruhancs/bank-go/db/sqlc"
	"github.com/ruhancs/bank-go/pb"
	"github.com/ruhancs/bank-go/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
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