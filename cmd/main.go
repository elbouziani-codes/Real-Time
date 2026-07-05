package main

import (
	"fmt"
	"log"
	"net/http"

	"realTime/config"
	"realTime/database"
	"realTime/handler"
)

func main() {
	con := config.Load()
	db, err := database.InitDB(con.DBPath)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println("Database Initialised")
	defer db.Close()

	mux1 := http.NewServeMux()

	mux1.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Server ONE on port "+con.PortMux1[1:]+"!")
	})

	go ServeFrontend(con.PortMux2)
	handler.RegisterHandlers(mux1)
	// 4. Start Server Two in the main goroutine to block and keep the application alive
	log.Println("Starting Server Two on " + con.PortMux1 + "...")
	if err := http.ListenAndServe(con.PortMux1, mux1); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server Two failed: %v", err)
	}
}
