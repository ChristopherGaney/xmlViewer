package config

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

const (
  host     = "localhost"
  port     = 5432
  user     = "postgres"
  password = "parset"
  dbname   = "parserdata"
)

var DB *sql.DB

func InitDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
        "password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
      panic(err)
    }
    defer db.Close()
    err = db.Ping()
    if err != nil {
      panic(err)
    }
}