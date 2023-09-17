package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" //driver postgres
	"github.com/ruhancs/bank-go/api"
	db "github.com/ruhancs/bank-go/db/sqlc"
	"github.com/ruhancs/bank-go/grpcapi"
	"github.com/ruhancs/bank-go/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	serverAddress = "localhost:8000"
	grpcAddress = "localhost:8001"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env")
	}

	dbDriver := os.Getenv("DB_DRIVER")
	dbSource := os.Getenv("DB_SOURCE")

	conn,err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db")
	}

	store := db.Newstore(conn)
	go runGinServer(store)
	runGRPCServer(store)
}

func runGRPCServer(store db.Store) {
	//criar novo grpc-server
	server,err := grpcapi.NewServer(store)
	if err != nil {
		log.Fatal("cannot instance server: ",err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterBankServer(grpcServer, server)
	//verificar os grpc disponiveis no servidor
	reflection.Register(grpcServer)

	listener,err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatal("cannot create listener: ",err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server: ")
	}
}

func runGinServer(store db.Store) {
	server,err := api.NewServer(store)
	if err != nil {
		log.Fatal("cannot instance server: ",err)
	}
	
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server")
	}

}