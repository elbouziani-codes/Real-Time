package main


import (
	"fmt"
	"context"
	"realTime/internal/repository"
	"realTime/internal/database"
	"realTime/internal/config"

)


func main() {
	conf := config.Load()
	id := "cc9bcbfc-104f-441e-b02d-8874c7cb320e"	
	 db, err := database.Open(conf.DBPath)
	if err != nil {
			fmt.Println(err)
			return
	}
	// new repo
	userRepo := repository.NewUserRepo(db)

	// new service
	user, err := userRepo.GetByID(context.Background(), id)
	if err != nil {
		fmt.Println(err, "tes")
		return
	}
	fmt.Println(user)
}
