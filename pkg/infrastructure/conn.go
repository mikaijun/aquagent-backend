package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

var Conn *sql.DB

func init() {
	var err error
	var envFile string
	goEnv := os.Getenv("GO_ENV")

	if goEnv == "" {
		envFile = ".env"
	} else {
		envFile = fmt.Sprintf(".env.%s", goEnv)
	}

	err = godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}

	user := os.Getenv("DBUser")
	password := os.Getenv("DBPassword")
	host := os.Getenv("DBHost")
	port := os.Getenv("DBPort")
	database := os.Getenv("DBName")

	Conn, err = sql.Open("postgres",
		fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, database))

	if err != nil {
		log.Fatal("OpenError: ", err)
	}

	if err = Conn.Ping(); err != nil {
		log.Fatal("PingError: ", err)
		log.Fatal("user: ", user)
		log.Fatal("password: ", password)
		log.Fatal("host: ", host)
		log.Fatal("port: ", port)
		log.Fatal("database: ", database)
	}
}
