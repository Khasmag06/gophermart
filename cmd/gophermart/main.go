package main

import (
	"github.com/Khasmag06/gophermart/config"
	"github.com/Khasmag06/gophermart/internal/client"
	"github.com/Khasmag06/gophermart/internal/handlers"
	"github.com/Khasmag06/gophermart/internal/repository"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := repository.NewDB(cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("unable to create database storage: %v", err)
	}
	accrualClient := client.NewAccrual(cfg.AccrualSystemAddress, repo)
	accrualClient.Run()
	s := handlers.NewService(accrualClient, repo)
	r := chi.NewRouter()

	r.Mount("/", s.Route())
	log.Fatal(http.ListenAndServe(cfg.RunAddress, r))
}
