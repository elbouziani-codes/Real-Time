package main

import (
	"context"
	"fmt"
	"realTime/internal/config"
	"realTime/internal/repository"
	"realTime/internal/service"
	"realTime/internal/repository/sqlite"
)

func main() {
	conf := config.Load()
	db, err := sqlite.Open(conf.DBPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	// new repo
	commentRepo := repository.NewPostRepo(db)
	commentService := service.NewPostService(commentRepo)	
	// new service
	comments, err := commentService.GetPosts(context.Background(), 10, 0)
	if err != nil {
		fmt.Println(err, "tes")
		return
	}
	fmt.Println(comments[0])
}
