package main

import (
	"log"
	"net/http"

	"RealTime/database"
	"RealTime/handler"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer db.Close()
	mux := http.NewServeMux()
	handler.Router(mux, db)
	log.Fatalln(http.ListenAndServe(":8080", mux))
}
