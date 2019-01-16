package main

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

var db *sql.DB

func InitDB() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
        "password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    db, err = sql.Open("postgres", psqlInfo)
    fmt.Println("initializing db")
    if err != nil {
      panic(err)
    }
    
    err = db.Ping()
    if err != nil {
      panic(err)
    }
}
