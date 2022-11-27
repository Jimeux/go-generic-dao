package db

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Jimeux/go-generic-dao/config"
)

// ErrNotFound is returned when no row is found for an SQL query.
var ErrNotFound = errors.New("database record not found")

var database *sql.DB

func DB() *sql.DB {
	return database
}

func Init() {
	c := config.Instance().Database
	dbString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.User, c.Password, c.Host, c.Port, c.Name)

	db, err := sql.Open("mysql", dbString)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to open database connection: %w", err))
	}
	db.SetConnMaxLifetime(time.Minute * 10)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	if err := db.Ping(); err != nil {
		log.Fatal(fmt.Errorf("failed to ping database: %w", err))
	}
	database = db
}

func Close() {
	_ = database.Close()
}
