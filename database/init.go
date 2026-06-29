package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)
var DB *sql.DB
func InitDB()  error{
	var err error
	DB , err := sql.Open("sqlite3","./realTime.db")
	if err != nil{
		return err
	}
	err = readShema()
	
	if err != nil{
		return err
	}
	
	return DB.Ping()
}

func readShema() err{
	bytes , err := os.ReadFile("shema.sql")
	if err != nil {
		return err
	}
	DB.Exec("INSERT ")
	return err
}