package grpcapi

import (
	db "github.com/ruhancs/bank-go/db/sqlc"
	"github.com/ruhancs/bank-go/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func converterUser(user db.User) *pb.User {
	return &pb.User{
		Username: user.Username,
		Fullname: user.FullName,
		Email: user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}