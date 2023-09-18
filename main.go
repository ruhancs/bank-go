package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/golang-migrate/migrate/v4" //migracoes automaticas do db
	_ "github.com/golang-migrate/migrate/v4/source/file"//driver de migracao via local file
	_ "github.com/golang-migrate/migrate/v4/database/postgres"//driver de migracao do db
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

	runDBMigrations(os.Getenv("MIGRATION_URL"), dbSource)

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

//migratioURL pasta que contem os arquivos das migracoes
func runDBMigrations(migratioURL, dbSource string) {
	migration,err := migrate.New(migratioURL,dbSource)
	if err != nil {
		log.Fatal("cannot create migrate instance: ",err)
	}
	
	//verificar se ocorreu erro nas migracoes e se o erro é sem altercoes nas migracoes
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("cannot run migrations up: ",err)
	}

	log.Println("db migrations successfuly")
}