package main

import (
	"context"
	"log"

	"realTime/database"
	"realTime/domain"
	"realTime/repository"
)

func main() {
	db, err := database.Open("test.db")
	if err != nil {
		log.Fatal(err)
	}
	userRepo := repository.NewUserRepo(db)
	testCreateUser(userRepo)
}

func testCreateUser(userRepo *repository.UserRepo) {
	user := domain.User{
		NickName:  "HMAR",
		LastName:  "ZARHON",
		FirstName: "younes",
		Age:       5,
		Gender:    "famme",
		Email:     "test@test.com",
	}
	err := userRepo.CreateUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
}
