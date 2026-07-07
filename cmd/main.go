package main

import (
	"fmt"
	"log"

	"realTime/config"
	"realTime/database"
	"realTime/repository"
)

func main() {
	conf := config.Load()
	db, err := db.Open(conf.DBPath)
	if err != nil {
		log.Fatalln(err)
		return
	}
	// new repo
	repo := repository.New(db)
	fmt.Println(repo)
	// new service
	// new handler 
	log.Println("Database Initialised")
}
