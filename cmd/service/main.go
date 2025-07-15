package main

import (
	"log"
	"net/http"

	"github.com/frederickmarvel/supernova/internal/config"
	"github.com/frederickmarvel/supernova/internal/db"
	"github.com/frederickmarvel/supernova/internal/router"
)

func main() {
	cfg := config.Load()
	database := db.New(cfg)
	r := router.New(database)

	log.Println("starting server on port 8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)
	}
}
