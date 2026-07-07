package main

import (
	"log"
	"net/http"
)

func ServeFrontend(port string) {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)
	log.Println("Starting Server Two on " + port[1:] + "...")

	if err := http.ListenAndServe(port, mux); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server Two failed: %v", err)
	}
}
