package main

import (
	"fmt"
	"log"
	"os"

	"realTime/config"
	"realTime/database"
	"realTime/repository"

	"github.com/projectdiscovery/interactsh/pkg/server"
)

func main() {
	conf := config.Load()

	db, err := db.Open(conf.DBPath)

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
