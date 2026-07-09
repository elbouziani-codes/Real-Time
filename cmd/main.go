package main

import (
	"fmt"
	"log"
	"os"

	"realTime/config"
	"realTime/database"
	"realTime/repository"

)

func main() {
	conf := config.Load()

	db, err := database.Open(conf.DBPath)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// new repo
	repo := repository.New(db)
	fmt.Println(repo)
	// new service
	// new handler 
	log.Println("Database Initialised")	
}
