package grpcapi

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	db "github.com/ruhancs/bank-go/db/sqlc"
	"github.com/ruhancs/bank-go/pb"
	"github.com/ruhancs/bank-go/util"
	"github.com/ruhancs/bank-go/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error loading .env")
	}

	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil,invalidArgumentError(violations)
	}

	user,err := server.store.GetUser(ctx,req.GetUsername())
	if err !=  nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "username does not registered")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}
	
	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credential")
	}

	tokenDuration,_ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_DURATION")) 
	duration := time.Duration(tokenDuration * int(time.Minute))
	accessToken,accessTokenPayload,err := server.tokenMaker.CreateToken(req.Username,duration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	}

	//gerar refresh token
	refreshTokenDuration,_ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DURATION")) 
	refreshDuration := time.Duration(refreshTokenDuration * int(time.Hour))
	refreshToken,refreshTokenPayload,err:= server.tokenMaker.CreateToken(user.Username,refreshDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate refresh token")
	}

	mtdt := server.extractMetadata(ctx)

	session,err := server.store.CreateSession(ctx,db.CreateSessionParams{
		ID: refreshTokenPayload.ID,
		Username: user.Username,
		RefreshToken: refreshToken,
		UserAgent: mtdt.UserAgent,
		ClientIp: mtdt.ClientIP,
		IsBlocked: false,
		ExpiresAt: refreshTokenPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate session")
	}

	resp := &pb.LoginUserResponse{
		User: converterUser(user),
		SessionId: session.ID.String(),
		AccessToken: accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshToken: refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
	}

	return resp,nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}