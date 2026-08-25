package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ishankgupta95/olx-api/internal/config"
	"github.com/ishankgupta95/olx-api/internal/db"
	"github.com/ishankgupta95/olx-api/internal/handlers"
)

func main() {

	cfg := config.MustLoad()
	db, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("main.db.Connect: %v", err)
	}
	fmt.Println("database connected")
	fmt.Println("Starting olx server...")

	lh := handlers.NewListingHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("server is listening on %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed %v", err)
	}
}
