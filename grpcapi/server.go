package grpcapi

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	db "github.com/ruhancs/bank-go/db/sqlc"
	"github.com/ruhancs/bank-go/pb"
	"github.com/ruhancs/bank-go/token"
	"github.com/ruhancs/bank-go/worker"
)

type Server struct {
	pb.UnimplementedBankServer//para nao precisar implementar as funcoes ao utilizar RegisterBankServer
	store db.Store
	tokenMaker token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(store db.Store, taskDistributor worker.TaskDistributor) (*Server,error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env")
	}
	symmetricKey := os.Getenv("TOKEN_SYMMETRIC_KEY")

	tokenMaker,err := token.NewPasetMaker(symmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store: store,
		tokenMaker: tokenMaker,
		taskDistributor: taskDistributor,
	}

	return  server,nil
}
