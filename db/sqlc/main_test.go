package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env")
	}

	dbDriver := os.Getenv("DB_DRIVER")
	dbSource := os.Getenv("DB_SOURCE")

	testDB,err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db")
	}

	//New vem db.go para inserir db conn em Queries
	testQueries = New(testDB)

	os.Exit(m.Run())
}