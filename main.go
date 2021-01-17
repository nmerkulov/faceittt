package main

import (
	"context"
	"faceittt/adapters/dummyemitter"
	"faceittt/adapters/postgres"
	"faceittt/adapters/web"
	"faceittt/application"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"os"
)

const (
	pgCSEnv = "PG_CONNSTRING"
)

func main() {
	repo, err := postgres.NewPGUserRepo(os.Getenv(pgCSEnv))
	if err != nil {
		log.Println(fmt.Errorf("main: %w", err))
		os.Exit(1)
	}
	repo.MustMigrate(context.Background())
	finder, err := postgres.NewPGUserFinder(os.Getenv(pgCSEnv))
	if err != nil {
		log.Println(fmt.Errorf("main: %w", err))
		os.Exit(1)
	}
	handlers := web.Router(web.WebParams{
		CreateUser: application.CreateUser(repo, dummyemitter.LogUserEvent),
		UpdateUser: application.UpdateUser(repo, finder, dummyemitter.LogUserEvent),
		DeleteUser: application.DeleteUser(repo, dummyemitter.LogUserEvent),
		FindUsers:  application.FindUsers(finder),
		FindUser:   application.FindUser(finder),
	})
	fmt.Println("listening :3000")
	if err := http.ListenAndServe(":3000", handlers); err != nil {
		log.Println(fmt.Errorf("main: %w", err))
		os.Exit(1)
	}
}
