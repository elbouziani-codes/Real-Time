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
	postRepo := repository.NewPostRepo(db)
	postService := service.NewPostService(postRepo)	
	// new service
	posts, err := postService.GetPosts(context.Background(), 10, 0)
	if err != nil {
		fmt.Println(err, "tes")
		return
	}
	fmt.Println(posts[0])
}
