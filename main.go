package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ruhancs/bank-go/api"
	db "github.com/ruhancs/bank-go/db/sqlc"
	_ "github.com/lib/pq" //driver postgres
)

const (
	serverAddress = "localhost:8000"
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
	server,err := api.NewServer(store)
	if err != nil {
		log.Fatal("cannot instance server: ",err)
	}

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server")
	}
}