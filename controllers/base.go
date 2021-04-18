package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nebisin/gopress/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Handler struct {
	Router *mux.Router
	DB *gorm.DB
}

func (handler *Handler) Initialize() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}


	handler.DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	} else {
		fmt.Println("Database connection is successful")
	}

	// Migrate the schema
	if err := handler.DB.AutoMigrate(&models.Post{}, &models.User{}); err != nil {
		log.Fatalf("Error auto migration: %v", err)
	}

	handler.Router = mux.NewRouter()

	handler.initializeRoutes()
}

func (handler *Handler) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, handler.Router))
}