package main

import (
	"log"

	"realTime/database"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println("Database Initialised")
	defer db.Close()
}
